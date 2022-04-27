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

	// playload
	Playload []byte `json:"playload"`

	// FormData
	FormData FormData `json:"form_data"`

	// header *http.Header
	Header *http.Header `json:"header"`
}
