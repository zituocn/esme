package esme

import (
	"reflect"
	"runtime"
	"strconv"
)

func Str2Int64(s string) int64 {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i64
}

func GetFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
