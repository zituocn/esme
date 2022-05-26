/*
task.go
http 任务

*/

package esme

import (
	"net/http"
)

// Task http Task
type Task struct {

	// Url 请求地址
	Url string `json:"url"`

	// Method 请求方法
	Method string `json:"method"`

	// Payload
	Payload []byte `json:"payload"`

	// FormData
	FormData FormData `json:"form_data"`

	// header *http.Header
	Header *http.Header `json:"header"`

	// Data 上下文数据传递
	Data map[string]interface{}
}
