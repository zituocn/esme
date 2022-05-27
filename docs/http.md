# http 请求参数和响应处理


### 一个简单的HTTP请求


*demo code*

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/zituocn/esme"
)

func main() {

	// 自定义header
	header := &http.Header{}
	header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
	header.Set("Content-Type", "application/json")

	// 设置payload
	body := `{"username":"testname","password":"1234567890"}`
	payLoad := []byte(body)

	// 请求一个天气预报的接口
	ctx := esme.HttpPost("https://tenapi.cn/wether/?city=%E6%88%90%E9%83%BD", header, payLoad)

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

	// 设置http代理
	ctx.SetProxy("http://10.10.10.10:8888")

	// 执行请求
	ctx.Do()
}

```

#### 设置FormData

```go

	// 设置formData
	formData := make(esme.FormData)
	formData["username"] = "testname"
	formData["password"] = "1234567890"

	// 请求一个天气预报的接口
	ctx := esme.HttpPost("https://tenapi.cn/wether/?city=%E6%88%90%E9%83%BD", header, formData)
```


#### 设置Cookie


一个cookie值

```go
	cookie := &http.Cookie{
		Name:  "token",
		Value: "cookie value",
	}
	ctx := esme.HttpPost("https://tenapi.cn/wether/?city=%E6%88%90%E9%83%BD", header, cookie)

```

多个cookie值 

```go
	cookie := make([]*esme.Cookie, 0)
	cookie = append(cookie, &esme.Cookie{
		Name:  "token",
		Value: "cookie value",
	})

	// 请求一个天气预报的接口
	ctx := esme.HttpPost("https://tenapi.cn/wether/?city=%E6%88%90%E9%83%BD", header, cookie)

```

```
此文档正在建设中...





```

### 更多文档


1. [http请求的参数设置&&响应处理](./docs/http.md)
2. [使用 `redis` 任务队列](./docs/job.md)
3. [在任务队列中，使用 `代理IP池` (多个代理IP使用)](./docs/proxy.md)
4. [和 `goquery`库的配合使用](./docs/html.md)
5. [和 `gjson` 库的配合使用](./doc/gjson.md)
6. [把数据存储到 `mysql` 中](./docs/db.md)