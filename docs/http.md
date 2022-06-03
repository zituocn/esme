# http 请求参数和响应处理


### HTTP 请求


*main.go*

```go
package main

import (
	"fmt"
	"net/http"
	"github.com/zituocn/esme"
)

func main() {

	// 请求一个天气预报的接口
	ctx := esme.HttpPost("https://tenapi.cn/wether/?city=%E6%88%90%E9%83%BD")

	// 成功的回调
	ctx.SetSucceedFunc(func(c *esme.Context) {
		fmt.Println("请求成功了:")
		fmt.Println("返回值 :", c.ToString())
	})
	// 失败的回调
	ctx.SetFailedFunc(func(c *esme.Context) {
		fmt.Println("请求出错了...")
		fmt.Println("返回状态值 :", c.Response.StatusCode)
		fmt.Println("返回值 :", c.ToString())
	})

	// 执行请求
	ctx.Do()
}

```

#### 设置 header

```go
header := &http.Header{}
header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
header.Set("Content-Type", "application/json")
ctx := esme.HttpPost("https://tenapi.cn/wether/?city=%E6%88%90%E9%83%BD", header)
```

#### 设置payload

```go
body := `{"username":"testname","password":"1234567890"}`
payLoad := []byte(body)
ctx := esme.HttpPost("https://tenapi.cn/wether/?city=%E6%88%90%E9%83%BD", payLoad)
```

#### 设置FormData

```go

formData := make(esme.FormData)
formData["username"] = "testname"
formData["password"] = "1234567890"
ctx := esme.HttpPost("https://tenapi.cn/wether/?city=%E6%88%90%E9%83%BD", formData)
```


#### 设置Cookie


一个cookie值

```go
cookie := &http.Cookie{
	Name:  "name",
	Value: "value",
}
ctx := esme.HttpPost("https://tenapi.cn/wether/?city=%E6%88%90%E9%83%BD", cookie)
```

多个cookie值 

```go
cookie := make([]*esme.Cookie, 0)
cookie = append(cookie, &esme.Cookie{
	Name:  "name",
	Value: "value",
})

ctx := esme.HttpPost("https://tenapi.cn/wether/?city=%E6%88%90%E9%83%BD", cookie)

```

#### 设置http代理

```go
ctx.SetProxy("http://10.10.10.10:8888")
```

#### 使用回调函数

请求开始的回调

```go
func (c *Context) SetStartFunc(fn CallbackFunc) *Context 

```

请求完成的回调

```go
func (c *Context) SetCompleteFunc(fn CallbackFunc) *Context
```

请求成功的回调

```go
func (c *Context) SetSucceedFunc(fn CallbackFunc) *Context 
```

请求失败的回调

```go

func (c *Context) SetFailedFunc(fn CallbackFunc) *Context 
```

需要重试的回调

```go
func (c *Context) SetRetryFunc(fn CallbackFunc) *Context

```

#### 设置http.client的transport

```go
func (c *Context) SetTransport(f func() *http.Transport) *Context 

```

演示代码

```go
ctx := HttpPost("https://tenapi.cn/wether/?city=%E6%88%90%E9%83%BD")

ctx.SetTransport(func() *http.Transport {
	return &http.Transport{
		MaxIdleConns:    100,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
})
```

---

### 响应处理

TODO:

---

### 更多文档


1. [http请求的参数设置&&响应处理](./docs/http.md)
2. [使用 `redis` 任务队列](./docs/job.md)
3. [在任务队列中，使用 `代理IP池` (多个代理IP使用)](./docs/proxy.md)
4. [和 `goquery`库的配合使用](./docs/html.md)
5. [和 `gjson` 库的配合使用](./doc/gjson.md)
6. [把数据存储到 `mysql` 中](./docs/db.md)