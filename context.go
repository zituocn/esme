/*
context.go
sam
2022-04-25
*/

package esme

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/zituocn/esme/logx"
)

// CallbackFunc call back func
type CallbackFunc func(*Context)

// Context request and response context
type Context struct {

	// http client
	client *http.Client

	// http request
	Request *http.Request

	// http response
	Response *http.Response

	// error
	Err error

	// callback for request start
	startFunc CallbackFunc

	// callback for successful request
	succeedFunc CallbackFunc

	// callback for failed request
	failedFunc CallbackFunc

	// callback to request try
	retryFunc CallbackFunc

	// callback for request completion
	completeFunc CallbackFunc

	// reqeusted task
	Task *Task

	// RespBody []byte returned by the request
	RespBody []byte

	// Param context parameter
	Data map[string]interface{}

	// sleepTime sheep time
	sleepTime time.Duration

	// isDebug debug mode switch
	isDebug bool

	// execution time
	execTime time.Duration
}

// Do execute current request
func (c *Context) Do() {
	var (
		bodyBytes []byte
		err       error
	)

	// set request body
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// sheep
	if c.sleepTime > 0 {
		time.Sleep(c.sleepTime)
	}

	// callback to execute start request
	if c.startFunc != nil {
		c.startFunc(c)
	}

	// start time
	startTime := time.Now()

	// start executing the request
	c.Response, c.Err = c.client.Do(c.Request)

	if c.Err != nil {
		//context deadline exceeded retry
		if c.retryFunc != nil {
			logx.Warnf("[%s] callback -> %s", "deadline", GetFuncName(c.retryFunc))
			c.retryFunc(c)
			return
		} else {
			logx.Errorf("request error: %s", c.Err.Error())
			return
		}
	}

	defer func(c *Context) {
		if c.Response != nil {
			err = c.Response.Body.Close()
			if err != nil {
				logx.Errorf("response body close error: %s", err.Error())
			}
		}
	}(c)

	c.execTime = time.Now().Sub(startTime)

	// gzip decode
	if c.Response.Header.Get("Content-Encoding") == "gzip" {
		c.Response.Body, err = gzip.NewReader(c.Response.Body)
		if err != nil {
			logx.Errorf("unzip failed: %s", err.Error())
			return
		}
	}

	// http response
	if c.Response != nil {
		code := c.Response.StatusCode
		status := GetStatusCodeString(code)
		body, err := ioutil.ReadAll(c.Response.Body)
		if err != nil {
			logx.Errorf("read response body error: %s", err.Error())
			logx.Debugf("task: %v", c.Task)
			return
		}
		c.RespBody = body
		switch status {
		case "success":
			// callback success function
			if c.succeedFunc != nil {
				logx.Infof("[%s] callback -> %s", status, GetFuncName(c.succeedFunc))
				c.succeedFunc(c)
			}
		case "retry":
			// callback retry function
			if c.retryFunc != nil {
				logx.Warnf("[%s] callback -> %s", status, GetFuncName(c.retryFunc))
				c.retryFunc(c)
				c.Do()
			}
		case "fail":
			// callback failed function
			if c.failedFunc != nil {
				logx.Errorf("[%s] callback -> %s", status, GetFuncName(c.failedFunc))
				c.failedFunc(c)
			}
		default:
			logx.Warnf("Unhandled status code: %d", code)
		}

	}

	// callback completion function
	if c.completeFunc != nil {
		c.completeFunc(c)
	}

	// isDebug print
	if c.isDebug {
		c.debugPrint()
	}
}

// SetIsDebug set debug
func (c *Context) SetIsDebug(isDebug bool) *Context {
	c.isDebug = isDebug
	return c
}

// SetSleepTime Set sleep time for http requests
func (c *Context) SetSleepTime(sleepTime int) *Context {
	c.sleepTime = time.Duration(sleepTime * int(time.Millisecond))
	return c
}

// SetTimeOut Set the timeout for http requests
//	milli second 毫秒
func (c *Context) SetTimeOut(timeout int) *Context {
	c.client.Timeout = time.Duration(timeout * int(time.Millisecond))
	return c
}

// SetStartFunc Set the callback for the start of the request
func (c *Context) SetStartFunc(fn CallbackFunc) *Context {
	c.startFunc = fn
	return c
}

// SetCompleteFunc Set the callback for request completion
func (c *Context) SetCompleteFunc(fn CallbackFunc) *Context {
	c.completeFunc = fn
	return c
}

// SetSucceedFunc Callback after successful setting
func (c *Context) SetSucceedFunc(fn CallbackFunc) *Context {
	c.succeedFunc = fn
	return c
}

// SetFailedFunc Callback after setting failure
func (c *Context) SetFailedFunc(fn CallbackFunc) *Context {
	c.failedFunc = fn
	return c
}

// SetRetryFunc Set callback for retry
func (c *Context) SetRetryFunc(fn CallbackFunc) *Context {
	c.retryFunc = fn
	return c
}

// SetProxy set http proxy
func (c *Context) SetProxy(httpProxy string) *Context {
	if httpProxy == "" {
		return c
	}
	proxy, _ := url.Parse(httpProxy)
	transport := getDefaultTransport()
	transport.Proxy = http.ProxyURL(proxy)
	c.client.Transport = transport
	return c
}

// SetTransport Set the client's Transport
func (c *Context) SetTransport(f func() *http.Transport) *Context {
	c.client.Transport = f()
	return c
}

// SetProxyLib set proxy lib
func (c *Context) SetProxyLib(lib *ProxyLib) *Context {
	if lib == nil {
		return c
	}
	ip, _ := lib.Get()
	c.SetProxy(ip)
	return c
}

// ToByte response body to []byte
func (c *Context) ToByte() []byte {
	if c.RespBody != nil {
		return c.RespBody
	}
	return []byte("")
}

// ToString response body to string
func (c *Context) ToString() string {
	if c.RespBody != nil {
		return string(c.RespBody)
	}
	return ""
}

// ToSection get json string by path
//		use gjson.Get func
//		like: ctx.ToSection("body.data")
func (c *Context) ToSection(path string) string {
	s := c.ToString()
	if s != "" {
		return gjson.Get(s, path).String()
	}
	return ""
}

// ToJSON returns error
//	response body to struct or map or slice
func (c *Context) ToJSON(v interface{}) error {
	if c.RespBody != nil {
		return json.Unmarshal(c.RespBody, &v)
	}
	return errors.New("response body is nil")
}

// ToHTML returns string
//	response body to html code
func (c *Context) ToHTML() string {
	s := c.ToString()
	return strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&#34;", `"`,
		"&#39;", "'",
	).Replace(s)
}

// GetExecTime get request execution time
func (c *Context) GetExecTime() time.Duration {
	return c.execTime
}

/*
private
*/

// debugPrint print request and response detail
func (c *Context) debugPrint() {

	fmt.Printf("%s %v \n", leftText("URL:"), c.Request.URL)
	fmt.Printf("%s %v \n", leftText("Method:"), c.Request.Method)
	fmt.Printf("%s %v \n", leftText("Request Header:"), c.Request.Header)
	fmt.Printf("%s %v \n", leftText("Response code:"), c.Response.StatusCode)
	fmt.Printf("%s %v \n", leftText("Response Header:"), c.Response.Header)

}

func (c *Context) reset() {
}

func leftText(s string) string {
	return fmt.Sprintf("%15s", s)
}
