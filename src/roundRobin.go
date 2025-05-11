package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

var current uint32

func (c *core) getNextBackend() string {
	index := atomic.AddUint32(&current, 1)
	return c.urls[(int(index)-1)%len(c.urls)]
}

func (c *core) handlerRoundRobin(w http.ResponseWriter, r *http.Request) {
	targetURL, err := url.Parse(c.getNextBackend())
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	r.Host = targetURL.Host
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		http.Error(rw, err.Error(), http.StatusBadGateway)
	}
	proxy.ServeHTTP(w, r)
}
