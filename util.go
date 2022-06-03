package esme

import (
	"math/rand"
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
