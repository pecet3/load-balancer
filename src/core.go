package main

import (
	"log"
	"net/http"
	"strconv"
)

type core struct {
	cfg               *Config
	currentSrvIndex   int
	currentWsSrvIndex int
	isWsSrv           bool
	db                *DB
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
	isWsSrv := false
	for _, srv := range cfg.Servers {
		if srv.IsWsCandidate {
			isWsSrv = true
			break
		}
	}
	return &core{
		cfg:     cfg,
		db:      db,
		isWsSrv: isWsSrv,
	}
}

func (c *core) RunLoadBalancer() {
	go c.getCurrentSrvIndexLoop()

	if c.isWsSrv {
		if c.cfg.IsDbLogging {
			http.HandleFunc("/", c.handlerWsWithLogs)
		} else {
			http.HandleFunc("/", c.handlerWs)
		}
	} else {
		if c.cfg.IsDbLogging {
			http.HandleFunc("/", c.handlerWithLogs)
		} else {
			http.HandleFunc("/", c.handler)
		}
	}

	log.Println("Listening on: ", c.cfg.Port)
	addr := ":" + strconv.Itoa(c.cfg.Port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalln("Server error: ", err)
	}
}
