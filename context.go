package esme

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"github.com/zituocn/gow/lib/logy"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

/*
context.go
*/

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
	StartFunc CallbackFunc

	// 成功的回调
	SucceedFunc CallbackFunc

	// 请求失败的回调
	FailedFunc CallbackFunc

	// 重试的回调
	RetryFunc CallbackFunc

	// 请求完成的回调
	CompleteFunc CallbackFunc

	// 请求的任务
	Task *Task

	// 请求返回的[]byte
	RespBody []byte

	Param map[string]interface{}
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
			if c.SucceedFunc != nil {
				logy.Infof("[%s] callback -> %s", status, GetFuncName(c.SucceedFunc))
				c.SucceedFunc(c)
			}
		case "retry":
			if c.RetryFunc != nil {
				logy.Warnf("[%s] callback -> %s", status, GetFuncName(c.RetryFunc))
				c.RetryFunc(c)
				c.Do()
			}
		case "fail":
			if c.FailedFunc != nil {
				logy.Errorf("[%s] callback -> %s", status, GetFuncName(c.FailedFunc))
				c.FailedFunc(c)
			}
		default:
			logy.Warnf("Unhandled status code: %d", code)
		}

	}

}

// SetSucceedFunc 设置成功后的回调
func (c *Context) SetSucceedFunc(fn CallbackFunc) *Context {
	c.SucceedFunc = fn
	return c
}

// SetFailedFunc 设置失败后的回调
func (c *Context) SetFailedFunc(fn CallbackFunc) *Context {
	c.FailedFunc = fn
	return c
}

// SetRetryFunc 设置重试的回调
func (c *Context) SetRetryFunc(fn CallbackFunc) *Context {
	c.RetryFunc = fn
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

func (c *Context) reset() {

}
