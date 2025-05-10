package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

var backendURLs = []string{
	"http://localhost:8083",
	"http://localhost:8082",
}

var current uint64

func getNextBackend() string {
	index := atomic.AddUint64(&current, 1)
	return backendURLs[(int(index)-1)%len(backendURLs)]
}

func handler(w http.ResponseWriter, r *http.Request) {
	targetURL, err := url.Parse(getNextBackend())
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
