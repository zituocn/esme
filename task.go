/*
task.go
http task
2022-04-25
*/

package esme

import (
	"net/http"
)

// Task http Task
type Task struct {

	// Url request address
	Url string `json:"url"`

	// Method request method
	Method string `json:"method"`

	// Payload request payload
	Payload []byte `json:"payload"`

	// FormData request formData
	FormData FormData `json:"form_data"`

	// header *http.Header
	Header *http.Header `json:"header"`

	// Data Contextual data passing
	Data map[string]interface{}
}
