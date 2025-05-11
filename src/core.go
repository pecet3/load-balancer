package main

import (
	"log"
	"net/http"
	"strconv"
)

type core struct {
	cfg        *Config
	urls       []string
	statusURLs []string
}

func NewCore(cfg *Config) *core {
	var urls []string
	for _, srv := range cfg.Servers {
		urls = append(urls, srv.URL)
	}
	var statusURLs []string
	for _, srv := range cfg.Servers {
		statusURLs = append(statusURLs, srv.StatusURL)
	}
	return &core{
		cfg:        cfg,
		urls:       urls,
		statusURLs: statusURLs,
	}
}

func (c *core) RunLoadBalancer() {
	log.Println("Listening on: ", c.cfg.Port)

	if c.cfg.IsRoundRobin {
		http.HandleFunc("/", c.handlerRoundRobin)
	}

	addr := ":" + strconv.Itoa(c.cfg.Port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("Błąd serwera: ", err)
	}
}
