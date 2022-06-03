/*
proxy_lib.go
代理ip库操作封装
*/

package esme

import (
	"sync"
	"sync/atomic"
)

// ProxyLib http proxy lib
type ProxyLib struct {
	mux *sync.Mutex
	ips []string
	num int32
	len int
}

// NewProxyLib return  new ProxyLib
func NewProxyLib() *ProxyLib {
	return &ProxyLib{
		ips: make([]string, 0),
		num: 0,
		len: 0,
		mux: &sync.Mutex{},
	}
}

// Add  proxyIP to ProxyLib
func (p *ProxyLib) Add(proxyIP *ProxyIP) {
	p.ips = append(p.ips, proxyIP.String())
	p.len++
}

// Del delete a ip by n
func (p *ProxyLib) Del(n int) {
	if (n + 1) > p.len {
		return
	}
	p.ips = append(p.ips[:n], p.ips[n+1:]...)
	p.len--
}

// Get get a ip
func (p *ProxyLib) Get() (string, int32) {
	if p.len == 0 {
		return "", 1
	}
	p.mux.Lock()
	// if end remove to first
	if p.num >= int32(p.len) {
		p.num = 0
	}
	ip := p.ips[p.num]
	p.mux.Unlock()
	n := p.num

	atomic.AddInt32(&p.num, 1)

	return ip, n
}
