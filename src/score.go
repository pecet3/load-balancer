package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Status struct {
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
}

func (c *core) handler(w http.ResponseWriter, r *http.Request) {
	c.logChan <- r
	targetURL, err := url.Parse(c.cfg.Servers[c.currentSrvIndex].URL)
	if err != nil {
		log.Println(getErr(r, err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	r.Host = targetURL.Host
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		log.Println(getErr(r, err))
		http.Error(rw, err.Error(), http.StatusBadGateway)
	}
	proxy.ServeHTTP(w, r)
}

// WEBSOCKETS

func (c *core) handlerWs(w http.ResponseWriter, r *http.Request) {
	c.logChan <- r
	srvUrl := c.cfg.Servers[c.currentSrvIndex].URL
	if isWebSocketRequest(r) {
		srvUrl = c.cfg.Servers[c.currentWsSrvIndex].URL
	}
	targetURL, err := url.Parse(srvUrl)
	if err != nil {
		log.Println(getErr(r, err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	r.Host = targetURL.Host
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		log.Println(getErr(r, err))
		http.Error(rw, err.Error(), http.StatusBadGateway)
	}
	proxy.ServeHTTP(w, r)
}

// loop

func (c *core) getCurrentSrvIndexLoop() {
	for {
		time.Sleep(time.Millisecond * time.Duration(c.cfg.StatusInterval))
		index := 0
		highestScore := 1e9
		indexWs := 0
		highestScoreWs := 1e9
		for i, srv := range c.cfg.Servers {
			status := &Status{}
			resp, err := http.Get(srv.StatusURL)
			if err != nil {
				log.Println(err)
				continue
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err := json.Unmarshal(body, status); err != nil {
				log.Println(err)
				continue
			}
			// calc score
			srvScore := (status.CPU * 0.8) + (status.Memory * 0.2)
			if srvScore < highestScore {
				highestScore = srvScore
				index = i
			}

			if srv.IsWsCandidate && srvScore < highestScoreWs {
				highestScoreWs = srvScore
				index = i
			}
			log.Println(srv.URL, "CPU:", status.CPU, "MEM:", status.Memory)
		}
		c.currentSrvIndex = index
		c.currentWsSrvIndex = indexWs
	}
}
