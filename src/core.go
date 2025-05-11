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
	score     float64
}

type core struct {
	cfg             *Config
	urls            []string
	servers         []*server
	currentSrvIndex uint32
	db              *DB
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
			score:     0,
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
	go c.getCurrentSrvIndexLoop()

	if c.cfg.IsDbLogging {
		http.HandleFunc("/", c.handlerScoreBasedWithLogs)
	} else {
		http.HandleFunc("/", c.handlerScoreBased)
	}

	log.Println("Listening on: ", c.cfg.Port)
	addr := ":" + strconv.Itoa(c.cfg.Port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalln("Server error: ", err)
	}
}
