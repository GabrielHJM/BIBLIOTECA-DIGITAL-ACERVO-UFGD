package httpclient

import (
	"net"
	"net/http"
	"time"
)

func NewOptimizedClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second, // Global timeout
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   10,
		},
	}
}
