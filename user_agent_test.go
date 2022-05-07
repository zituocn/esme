package esme

import (
	"fmt"
	"testing"
)

func Test_RandomUserAgent(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := RandomUserAgent(PCUserAgent)
		fmt.Println(s)
	}
}
