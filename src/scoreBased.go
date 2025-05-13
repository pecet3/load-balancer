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

func (c *core) handlerScoreBased(w http.ResponseWriter, r *http.Request) {
	targetURL, err := url.Parse(c.cfg.Servers[c.currentSrvIndex].URL)
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
func (c *core) handlerScoreBasedWithLogs(w http.ResponseWriter, r *http.Request) {
	go c.db.AddRequest(r)
	c.handlerScoreBased(w, r)
}

func (c *core) getCurrentSrvIndexLoop() {
	for {
		time.Sleep(time.Millisecond * time.Duration(c.cfg.StatusInterval))
		index := 0
		highestScore := 1e9
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
			if c.cfg.IsDbLogging {
				go c.db.AddStatus(srv.URL, status.CPU, status.Memory)
			}
		}
		c.currentSrvIndex = uint32(index)
	}
}
