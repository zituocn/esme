/*
context.go
*/

package esme

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/zituocn/esme/logx"
)

// CallbackFunc 回调函数
type CallbackFunc func(*Context)

// Context 上下文封装
type Context struct {

	// http client
	client *http.Client

	// http request
	Request *http.Request

	// http response
	Response *http.Response

	//ctx
	Ctx context.Context

	// error
	Err error

	// 请求开始的回调
	startFunc CallbackFunc

	// 成功的回调
	succeedFunc CallbackFunc

	// 请求失败的回调
	failedFunc CallbackFunc

	// 重试的回调
	retryFunc CallbackFunc

	// 请求完成的回调
	completeFunc CallbackFunc

	// 请求的任务
	Task *Task

	// 请求返回的[]byte
	RespBody []byte

	Param map[string]interface{}

	sleepTime time.Duration

	isDebug bool
}

// Do 执行当前请求
func (c *Context) Do() {
	var (
		bodyBytes []byte
		err       error
	)

	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// 休眠
	if c.sleepTime > 0 {
		time.Sleep(c.sleepTime)
	}

	// 执行开始请求的回调
	c.startFunc(c)

	// 开始执行请求
	c.Response, c.Err = c.client.Do(c.Request)
	if c.Err != nil {
		logx.Errorf("请求出错: %v", c.Err)
		return
	}

	defer func(c *Context) {
		if c.Response != nil {
			c.Response.Body.Close()
		}
	}(c)

	// gzip decode
	if c.Response.Header.Get("Content-Encoding") == "gzip" {
		c.Response.Body, err = gzip.NewReader(c.Response.Body)
		if err != nil {
			logx.Errorf("unzip failed :%v", err)
			return
		}
	}

	// http response
	if c.Response != nil {
		code := c.Response.StatusCode
		status := GetStatusCodeString(code)
		body, err := ioutil.ReadAll(c.Response.Body)
		if err != nil {
			logx.Errorf("read response body error : %v", err)
			return
		}
		c.RespBody = body
		switch status {
		case "success":
			// 回调成功函数
			if c.succeedFunc != nil {
				logx.Infof("[%s] callback -> %s", status, GetFuncName(c.succeedFunc))
				c.succeedFunc(c)
			}
		case "retry":
			if c.retryFunc != nil {
				logx.Warnf("[%s] callback -> %s", status, GetFuncName(c.retryFunc))
				c.retryFunc(c)
				c.Do()
			}
		case "fail":
			if c.failedFunc != nil {
				logx.Errorf("[%s] callback -> %s", status, GetFuncName(c.failedFunc))
				c.failedFunc(c)
			}
		default:
			logx.Warnf("Unhandled status code: %d", code)
		}

	}

	// isDebug print
	if c.isDebug {
		c.debugPrint()
	}

	// 执行请求完成的回调
	c.completeFunc(c)

}

// SetIsDebug 设置是否打印debug信息
func (c *Context) SetIsDebug(isDebug bool) *Context {
	c.isDebug = isDebug
	return c
}

// SetSleepTime 设置http请求的休眠时间
func (c *Context) SetSleepTime(sleepTime int) *Context {
	c.sleepTime = time.Duration(sleepTime * int(time.Millisecond))
	return c
}

// SetTimeOut 设置http请求的超时时间
//	milli second 毫秒
func (c *Context) SetTimeOut(timeout int) *Context {
	c.client.Timeout = time.Duration(timeout * int(time.Millisecond))
	return c
}

// SetStartFunc 设置请求开始的回调
func (c *Context) SetStartFunc(fn CallbackFunc) *Context {
	c.startFunc = fn
	return c
}

// SetCompleteFunc 设置请求完成的回调
func (c *Context) SetCompleteFunc(fn CallbackFunc) *Context {
	c.completeFunc = fn
	return c
}

// SetSucceedFunc 设置成功后的回调
func (c *Context) SetSucceedFunc(fn CallbackFunc) *Context {
	c.succeedFunc = fn
	return c
}

// SetFailedFunc 设置失败后的回调
func (c *Context) SetFailedFunc(fn CallbackFunc) *Context {
	c.failedFunc = fn
	return c
}

// SetRetryFunc 设置重试的回调
func (c *Context) SetRetryFunc(fn CallbackFunc) *Context {
	c.retryFunc = fn
	return c
}

// SetProxy 设置 http 代理
func (c *Context) SetProxy(httpProxy string) *Context {
	if httpProxy == "" {
		return c
	}
	proxy, _ := url.Parse(httpProxy)
	c.client.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}
	return c
}

// SetProxyFunc 设置代理方法
func (c *Context) SetProxyFunc(f func() *http.Transport) *Context {
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
