package esme

import (
	"fmt"
	"testing"
)

func Test_GetRandSleepTime(t *testing.T) {
	fmt.Println(GetRandSleepTime(1000, 3000))

	fmt.Println(GetRandSleepTime(1, 100))
	fmt.Println(GetRandSleepTime(0, -1))
}
