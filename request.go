package esme

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zituocn/gow/lib/logy"
	"golang.org/x/net/publicsuffix"
	"net"
	"net/http"
	"net/http/cookiejar"
	burl "net/url"
	"time"
)

/*
http_request.go
http请求的相关
*/

type Header map[string]string

type FormData map[string]string

func (h Header) haveObj() {
	if h == nil {
		h = Header{}
	}
}

func (h Header) Set(key, value string) Header {
	h.haveObj()
	h[key] = value
	return h
}

func (h Header) Delete(key string) Header {
	h.haveObj()
	delete(h, key)
	return h
}

// RequestTimeOut request 请求超时
//	毫秒
type RequestTimeOut int

type Cookie struct {
	Name     string
	Value    string
	HttpOnly bool
}

// HTTPGet 执行一个http get请求
func HTTPGet(url string, vs ...interface{}) *Context {
	return DoRequest(url, "GET", vs...)
}

// HTTPPost 执行一个http post请求
func HTTPPost(url string, vs ...interface{}) *Context {
	return DoRequest(url, "POST", vs...)
}

// HTTPPut 执行一个http put请求
func HTTPPut(url string, data []byte, vs ...interface{}) *Context {
	return DoRequest(url, "PUT", vs...)
}

// DoRequest 执行一个请求
//	不能执行回调
//	返回 Context
func DoRequest(url, method string, vs ...interface{}) *Context {
	ctx, err := NewRequest(url, method, vs...)
	if err != nil {
		logy.Errorf("DoRequest 错误 :%v", err)
		return nil
	}
	return ctx
}

// GetByte interface{} to []byte
func GetByte(obj interface{}) (data []byte) {
	if obj != nil {
		data, _ = json.Marshal(&obj)
	}
	return
}

// NewRequest returns esme.context and error
func NewRequest(url, method string, vs ...interface{}) (*Context, error) {
	u, errU := validUrl(url)
	if errU != nil {
		return nil, errU
	}
	req, err := http.NewRequest(method, u, nil)

	for _, v := range vs {
		switch vv := v.(type) {
		case FormData:
			formData := burl.Values{}
			for k, v := range vv {
				formData.Add(k, v)
			}
			req, err = http.NewRequest(method, u, bytes.NewBuffer([]byte(formData.Encode())))
			if err != nil {
				return nil, err
			}
		case []byte:
			{
				req, err = http.NewRequest(method, u, bytes.NewReader(vv))
				if err != nil {
					return nil, err
				}
			}
		default:

		}

	}
	req.Header = http.Header{}
	ctx := NewContext(req, vs...)
	return ctx, nil
}

// NewContext returns new Context
func NewContext(req *http.Request, vs ...interface{}) *Context {

	var (
		client         *http.Client
		requestTimeOut RequestTimeOut
	)

	for _, v := range vs {
		switch vv := v.(type) {
		case http.Header:
			for key, values := range vv {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}
		case *http.Header:
			for key, values := range *vv {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}
		case Header:
			for key, value := range vv {
				req.Header.Add(key, value)
			}
		case *http.Client:
			client = vv
		case *http.Cookie:
			req.AddCookie(vv)
		case []Cookie:
			for _, cookie := range vv {
				req.AddCookie(&http.Cookie{
					Name:     cookie.Name,
					Value:    cookie.Value,
					HttpOnly: cookie.HttpOnly,
				})
			}
		case []*Cookie:
			for _, cookie := range vv {
				req.AddCookie(&http.Cookie{
					Name:     cookie.Name,
					Value:    cookie.Value,
					HttpOnly: cookie.HttpOnly,
				})
			}
		case RequestTimeOut:
			requestTimeOut = vv
		case FormData:
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
		}

	}

	if client == nil {
		client = getDefaultClient()
	}

	if requestTimeOut > 0 {
		client.Timeout = time.Duration(requestTimeOut) * time.Millisecond
	}

	// set transport
	client.Transport = getDefaultTransport()

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		client.Jar = jar
	}

	//decode gzip
	if length := req.Header.Get("Content-Length"); length != "" {
		req.ContentLength = Str2Int64(length)
	}

	return &Context{
		client:  client,
		Request: req,
		Param:   make(map[string]interface{}),
	}
}

/*
private
*/

func getDefaultClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

func getDefaultTransport() *http.Transport {
	return &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     true,
	}
}

func validUrl(urlStr string) (string, error) {

	length := len(urlStr)
	if length < 7 {
		return "", fmt.Errorf("错误的 url : %s", urlStr)
	}

	if urlStr[:7] == "http://" || urlStr[:8] == "https://" {
		return urlStr, nil
	}

	return "http://" + urlStr, nil
}
