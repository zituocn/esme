package esme

import (
	"fmt"
	"testing"
)

func Test_RequestAndResponse(t *testing.T) {
	ctx := HttpGet("https://tenapi.cn/douyinresou/")
	ctx.SetSucceedFunc(func(c *Context) {
		fmt.Println("数据获取成功了...")
	})
	ctx.Do()
	str := ctx.ToString()
	fmt.Println(str)
}
