package esme

import (
	"fmt"
	"testing"
)

func Test_GetRandSleepTime(t *testing.T) {
	fmt.Println(GetRandSleepTime(1000, 3000))

	fmt.Println(GetRandSleepTime(1, 100))
	fmt.Println(GetRandSleepTime(0, -1))
}

func Test_GetFuncName(t *testing.T) {
	fmt.Println(GetFuncName(func(c *Context) {

	}))

}

func Test_BuilderHeader(t *testing.T) {
	s := `
	Accept: application/json, text/plain, */*
	Accept-Encoding: gzip, deflate, br
	Accept-Language: zh-CN,zh;q=0.9
	Connection: keep-alive
	Host: api.zhiyuantong.com
	Origin: https://b.zhiyuantong.com
	Referer: https://b.zhiyuantong.com/
	sec-ch-ua: " Not A;Brand";v="99", "Chromium";v="102", "Google Chrome";v="102"
	sec-ch-ua-mobile: ?0
	sec-ch-ua-platform: "macOS"
	Sec-Fetch-Dest: empty
	Sec-Fetch-Mode: cors
	Sec-Fetch-Site: same-site
	User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36
	X-Router: /datapanel/college/home
	X-Token: eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiIsImZvbyI6ImppZWppZSJ9.eyJpc3MiOiJEaWdnZyIsImF1ZCI6Ind3dy56aGl5dWFudG9uZy5jb20iLCJqdGkiOiIxNjg4OCIsImlhdCI6MTY1NTc4MDgxOSwiZXhwIjoxNjU1ODI0MDE5LCJ1aWQiOiJ7XCJpZFwiOjcwMzEzLFwibGFzdGxvZ2luaXBcIjpcIjE3MS4yMjEuMTQ2LjEwNlwiLFwibGFzdGxvZ2luZGF0ZVwiOjE2NTU3ODA4MTksXCJzaWduXCI6XCJcIn0ifQ.oZZ6MYz-T9bBxxtnTdrlHEYW_Wy6473PR-ttckSQ194`
	header := BuilderHeader(s)

	fmt.Println(header)
}

func Test_BuilderFormData(t *testing.T) {
	s := `page=1&limit=15&id=&nick_name=&mobile=&source_type=-100&gid=-100&sub_id=-100&prov_id=-100&vip_goods_id=-100&user_type=-100&stime=0&etime=0`

	formData := BuilderFormData(s)

	fmt.Println(formData)
}
