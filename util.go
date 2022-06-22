package esme

import (
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Str2Int64 str convert to int64
func Str2Int64(s string) int64 {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i64
}

// GetFuncName return the function name by reflection
func GetFuncName(i interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	if name != "" {
		return name[strings.LastIndex(name, "/")+1:]
	}
	return ""
}

// GetRandSleepTime Generate random from min to max
//	including min and max
func GetRandSleepTime(min, max int) int {
	if min < 1 {
		min = 1
	}
	if max < 1 {
		max = 1
	}
	if min == max {
		return min
	}
	rand.Seed(time.Now().Unix())
	n := rand.Intn(max - min)
	n = n + min
	return n
}

// BuilderHeader convert string headers to http.Header
func BuilderHeader(s string) *http.Header {
	header := &http.Header{}
	if s == "" {
		return header
	}
	s = strings.ReplaceAll(s, "\t", "")
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		line = strings.TrimLeft(line, "")
		line = strings.TrimRight(line, "")
		if line != "" {
			kvs := strings.Split(line, ": ")
			if len(kvs) == 2 {
				header.Set(strings.TrimSpace(kvs[0]), strings.TrimSpace(kvs[1]))
			}
		}
	}
	return header
}

// BuilderFormData convert string FormData to esme.FormData
//  page=1&limit=15&id=&nick_name=&mobile=&source_type=-100 to FormData
func BuilderFormData(s string) FormData {
	formData := FormData{}
	if s == "" {
		return formData
	}
	lines := strings.Split(s, "&")
	for _, line := range lines {
		if line != "" {
			kvs := strings.Split(line, "=")
			if len(kvs) == 2 {
				formData.Set(kvs[0], kvs[1])
			}
		}
	}
	return formData
}
