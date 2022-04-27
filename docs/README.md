# esme

一个go实现的多任务、多线程的网络请求SDK


### 1. 特性

* 基本的网络请求实现
* 自定义请求参数
* http 代理的使用
* http 请求 debug 模式
* http 请求的回调，包括成功回调、失败回调和重试回调
* 支持自己设置 http.Transport
* 内存中的任务队列实现
* redis 任务队列实现

### 2. 安装

```shell
go get -u github.com/zituocn/esme
```

### 3. 简单的HTTP请求

一个 http get请求的演示代码

*demo code*

```go
package main

import (
	"fmt"
	"github.com/zituocn/esme"
)

func main() {

    // 请求一个天气预报的接口
	ctx := esme.HttpGet("https://tenapi.cn/wether/?city=%E6%88%90%E9%83%BD")

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

*返回值*

```shell
2022/04/26 19:42:57.625 [I] context.go:120: [success] callback -> main.main.func1
请求成功了:
返回值 : ....
```

---
### 4. 使用任务队列

使用内存任务队列的演示

*demo code*

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/zituocn/esme"
)

var (
	// queue 一个内存中的任务队列
	queue = esme.NewMemQueue()
)

// AddTask 添加任务
func AddTask() {

	city := []string{"北京", "上海", "成都", "深圳", "西安"}

	header := &http.Header{}
	header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
	for _, item := range city {
		queue.Add(&esme.Task{
			Url:    fmt.Sprintf("%s%s", "https://tenapi.cn/wether/?city=", item),
			Method: "GET",
			Header: header,
		})
	}
}

func main() {

	// 生成任务队列
	AddTask()

	// 设置 任务参数
	job := esme.NewJob("wether", 1, queue, esme.JobOptions{
		SucceedFunc: func(ctx *esme.Context) {
			fmt.Println("成功的回调")
			fmt.Println("返回信息 :", ctx.ToString())
		},
		FailedFunc: func(ctx *esme.Context) {
			fmt.Println("失败的回调")
			fmt.Println("返回状态 :", ctx.Response.StatusCode)
		},
	})

	// 执行
	job.Do()

}
```


### 5. 更多文档

1. [http请求的参数设置&&响应处理](./docs/http.md)
2. [使用 `redis` 任务队列](./docs/job.md)
3. [在任务队列中，使用 `代理IP池` (多个代理IP使用)](./docs/proxy.md)
4. [和 `goquery`库的配合使用](./docs/html.md)
5. [把数据存储到 `mysql` 中](./docs/db.md)

### 6. 感谢&&参考


* [gathertool](https://github.com/mangenotwork/gathertool)
