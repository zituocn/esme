/*
proxy.go
sam
2022-04-25
*/

package esme

import (
	"fmt"
)

// ProxyIP proxy ip struct
//	http proxy
type ProxyIP struct {
	IP    string
	Port  int
	User  string
	Pass  string
	IsTLS bool
}

// NewProxyIP return http proxy
func NewProxyIP(ip string, port int, user, pass string, isTls bool) *ProxyIP {
	return &ProxyIP{
		IP:    ip,
		Port:  port,
		User:  user,
		Pass:  pass,
		IsTLS: isTls,
	}
}

// String return http proxy string
func (p *ProxyIP) String() string {
	scheme := "http"
	if p.IsTLS {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s:%s@%s:%d", scheme, p.User, p.Pass, p.IP, p.Port)
}
