package main

import (
	"log"
	"net/http"
	"strconv"
)

type core struct {
	cfg             *Config
	urls            []string
	currentSrvIndex uint32
	db              *DB
}

func NewCore(cfg *Config) *core {
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
		cfg: cfg,
		db:  db,
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
