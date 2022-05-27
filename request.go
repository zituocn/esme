/*
request.go
http请求的相关
sam
*/

package esme

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	burl "net/url"
	"strings"
	"time"

	"github.com/zituocn/esme/logx"
	"golang.org/x/net/publicsuffix"
)

const (
	defaultContentType = "text/html; charset=utf-8"
	defaultUserAgent   = "Go-http-client/esme/1.0"
)

// Header esme.Header
type Header map[string]string

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

// Cookie esme.Cookie
type Cookie struct {
	Name     string
	Value    string
	HttpOnly bool
}

// FormData esme.FormData
type FormData map[string]string


// HttpGet 执行一个http get请求
func HttpGet(url string, vs ...interface{}) *Context {
	return DoRequest(url, "GET", vs...)
}

// HttpPost 执行一个http post请求
func HttpPost(url string, vs ...interface{}) *Context {
	return DoRequest(url, "POST", vs...)
}

// HttpPut 执行一个http put请求
func HttpPut(url string, data []byte, vs ...interface{}) *Context {
	return DoRequest(url, "PUT", vs...)
}

// DoRequest 执行一个请求
//	返回 Context
func DoRequest(url, method string, vs ...interface{}) *Context {
	ctx, err := NewRequest(url, method, vs...)
	if err != nil {
		logx.Errorf("DoRequest 错误 :%v", err)
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
			if len(vv) > 0 {
				formData := burl.Values{}
				for k, v := range vv {
					formData.Set(k, v)
				}
				req, err = http.NewRequest(method, u, strings.NewReader(formData.Encode()))
				if err != nil {
					return nil, err
				}
			}
		case []byte:
			{
				if len(vv) > 0 {
					req, err = http.NewRequest(method, u, bytes.NewReader(vv))
					if err != nil {
						return nil, err
					}
					req.ContentLength = int64(len(vv))
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
		client *http.Client
		task   *Task
	)
	for _, v := range vs {
		switch vv := v.(type) {
		case *Task:
			task = vv
		case http.Header:
			for key, values := range vv {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}
		case *http.Header:
			if vv != nil {
				for key, values := range *vv {
					for _, value := range values {
						req.Header.Add(key, value)
					}
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
		case FormData:
			if len(vv) > 0 {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}
		}
	}

	if client == nil {
		client = getDefaultClient()
	}

	// set transport
	client.Transport = getDefaultTransport()

	// cookie jar
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		client.Jar = jar
	}

	if length := req.Header.Get("Content-Length"); length != "" {
		req.ContentLength = Str2Int64(length)
	}

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", defaultContentType)
	}

	if req.Header.Get("User-Agent") == "" || req.Header.Get("User-Agent") == "Go-http-client/2.0" {
		req.Header.Set("User-Agent", defaultUserAgent)
	}

	return &Context{
		client:  client,
		Request: req,
		Task:    task,
		Data:    make(map[string]interface{}),
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
		MaxIdleConns:    100,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func validUrl(urlStr string) (string, error) {
	length := len(urlStr)
	if length < 7 {
		return "", fmt.Errorf("error request url : %s", urlStr)
	}

	if urlStr[:7] == "http://" || urlStr[:8] == "https://" {
		return urlStr, nil
	}

	return "http://" + urlStr, nil
}
