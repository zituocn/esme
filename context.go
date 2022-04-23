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

	"github.com/zituocn/gow/lib/logy"
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
	)

	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// 休眠
	if c.sleepTime > 0 {
		time.Sleep(c.sleepTime)
	}

	// 开始执行请求
	c.Response, c.Err = c.client.Do(c.Request)
	if c.Err != nil {
		logy.Errorf("请求出错: %v", c.Err)
		return
	}

	defer func(c *Context) {
		if c.Response != nil {
			c.Response.Body.Close()
		}
	}(c)

	// gzip decode
	if c.Response.Header.Get("Content-Encoding") == "gzip" {
		c.Response.Body, _ = gzip.NewReader(c.Response.Body)
	}

	// http response
	if c.Response != nil {
		code := c.Response.StatusCode
		status := GetStatusCodeString(code)
		switch status {
		case "success":
			body, err := ioutil.ReadAll(c.Response.Body)
			if err != nil {
				logy.Errorf("读取 response body error : %v", err)
				return
			}
			c.RespBody = body
			// 回调成功函数
			if c.succeedFunc != nil {
				logy.Infof("[%s] callback -> %s", status, GetFuncName(c.succeedFunc))
				c.succeedFunc(c)
			}
		case "retry":
			if c.retryFunc != nil {
				logy.Warnf("[%s] callback -> %s", status, GetFuncName(c.retryFunc))
				c.retryFunc(c)
				c.Do()
			}
		case "fail":
			if c.failedFunc != nil {
				logy.Errorf("[%s] callback -> %s", status, GetFuncName(c.failedFunc))
				c.failedFunc(c)
			}
		default:
			logy.Warnf("Unhandled status code: %d", code)
		}

	}

	if c.isDebug {
		c.debugPrint()
	}

}

// SetIsDebug set is debug
//	if isDebug is true , print http debug
func (c *Context) SetIsDebug(isDebug bool) *Context {
	c.isDebug = isDebug
	return c
}

// SetSleepTime set request sleep time
func (c *Context) SetSleepTime(sleepTime int) *Context {
	c.sleepTime = time.Duration(sleepTime * int(time.Millisecond))
	return c
}

// SetTimeOut http request timeout
//	milli second 毫秒
func (c *Context) SetTimeOut(timeout int) *Context {
	c.client.Timeout = time.Duration(timeout * int(time.Millisecond))
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

// SetProxy set http proxy
func (c *Context) SetProxy(httpProxy string) {
	proxy, _ := url.Parse(httpProxy)
	c.client.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}

}

// SetProxyFunc set transport func
func (c *Context) SetProxyFunc(f func() *http.Transport) {
	c.client.Transport = f()
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

// debugPrint print request and response detail
func (c *Context) debugPrint() {

	fmt.Println("method =", c.Request.Method)
	fmt.Println("url =", c.Request.URL)

}

func (c *Context) reset() {
}
