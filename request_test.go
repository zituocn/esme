package esme

import (
	"fmt"
	"testing"
)

func Test_RequestAndResponse(t *testing.T) {
	ctx := DoRequest("https://test.api.ymzy.cn/static_service/v1/allow/index/info", "GET", "application/json", nil)
	ctx.SetSucceedFunc(func(c *Context){
		fmt.Println("数据成功了.....")
	})
	str := ctx.ToString()
	fmt.Println(str)
}
