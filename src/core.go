package main

import (
	"log"
	"net/http"
	"strconv"
)

type server struct {
	url       string
	statusURL string
	isDown    bool
	cpuUsage  float64
}

type core struct {
	cfg     *Config
	urls    []string
	servers []*server
	db      *DB
}

func NewCore(cfg *Config) *core {
	var urls []string
	var servers []*server
	for _, srv := range cfg.Servers {
		urls = append(urls, srv.URL)
		servers = append(servers, &server{
			url:       srv.URL,
			statusURL: srv.StatusURL,
			isDown:    false,
			cpuUsage:  0,
		})
	}
	var db *DB
	if cfg.IsDbLogging {
		if ndb, err := NewDB(); err == nil {
			db = ndb
		} else {
			log.Fatalln("Error with initalize database: ", err)
		}
	} else {
		db = nil
	}
	return &core{
		cfg:     cfg,
		urls:    urls,
		servers: servers,
		db:      db,
	}
}

func (c *core) RunLoadBalancer() {
	log.Println("Listening on: ", c.cfg.Port)

	if c.cfg.IsRoundRobin {
		if c.cfg.IsDbLogging {
			http.HandleFunc("/", c.handlerRoundRobinWithLogs)
		} else {
			http.HandleFunc("/", c.handlerRoundRobin)
		}
	}

	addr := ":" + strconv.Itoa(c.cfg.Port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalln("Server error: ", err)
	}
}
