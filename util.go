package esme

import (
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"time"
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

// GetRandSleepTime 产生从min到max之间的随机机
//	包括 min和max
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
