package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type core struct {
	cfg               *Config
	currentSrvIndex   int
	currentWsSrvIndex int
	isWsSrv           bool
	logChan           chan (*http.Request)
}

func NewCore(cfg *Config) *core {
	isWsSrv := false
	for _, srv := range cfg.Servers {
		if srv.IsWsCandidate {
			isWsSrv = true
			break
		}
	}
	return &core{
		cfg:     cfg,
		isWsSrv: isWsSrv,
		logChan: make(chan *http.Request, cfg.LogBuffSize),
	}
}

func (c *core) RunLoadBalancer() {
	go c.getCurrentSrvIndexLoop()
	go c.runLogger()
	if c.isWsSrv {
		http.HandleFunc("/", c.handlerWs)
	} else {
		http.HandleFunc("/", c.handler)
	}

	log.Println("Listening on:", c.cfg.Port)
	addr := ":" + strconv.Itoa(c.cfg.Port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalln("Server error: ", err)
	}
}

func (c *core) runLogger() {
	for {
		req := <-c.logChan
		log.Println(getInfo(req))
	}
}

func getInfo(req *http.Request) string {
	return fmt.Sprint("| ", req.URL.RequestURI(), " | IP:", getIP(req), " | WS:", isWebSocketRequest(req))
}

func getErr(req *http.Request, err error) string {
	return fmt.Sprint(getInfo(req), " | ERROR:", err.Error())
}
