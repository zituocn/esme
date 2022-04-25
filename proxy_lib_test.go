package esme

import (
	"fmt"
	"testing"
)

var (
	proxyLib = NewProxyLib()
)

func Add() {
	for i := 0; i < 100; i++ {
		proxyLib.Add(
			&ProxyIP{
				IP:    fmt.Sprintf("10.10.10.%d", i+1),
				Port:  8888,
				User:  "",
				Pass:  "",
				IsTLS: false,
			})
	}
}

func Test_ProxyLibGet(t *testing.T) {
	Add()
	for i := 0; i < 10; i++ {
		ip, n := proxyLib.Get()
		fmt.Println(ip, n)
	}
}
