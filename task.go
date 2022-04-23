/*
task.go
http 任务

*/

package esme

import (
	"net/http"
	"sync"
)

// Task http Task
type Task struct {

	// Url 请求地址
	Url string

	// 请求方法
	Method string

	// playload
	Playload []byte

	// FormData
	FormData FormData

	// header
	Header *http.Header

	once *sync.Once
}
