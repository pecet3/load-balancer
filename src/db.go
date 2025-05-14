package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	d *sql.DB
}

func NewDB() (*DB, error) {
	err := os.MkdirAll("data/database", os.ModePerm)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", "file:data/database/store.db?_journal_mode=WAL")
	if err != nil {
		return nil, err
	}

	createRequestsTable := `
	CREATE TABLE IF NOT EXISTS requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT NOT NULL,
		path TEXT NOT NULL,
		is_websocket BOOL NOT NULL,
		created_at DATETIME NOT NULL
	);
	`
	createStatusesTable := `
	CREATE TABLE IF NOT EXISTS statuses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		cpu REAL NOT NULL,
		memory REAL NOT NULL,
		created_at DATETIME NOT NULL
	);
	`

	if _, err := db.Exec(createRequestsTable); err != nil {
		return nil, err
	}
	if _, err := db.Exec(createStatusesTable); err != nil {
		return nil, err
	}

	return &DB{d: db}, nil
}
func (db DB) AddStatus(path string, cpu float64, memory float64) {
	insertQuery := `
	INSERT INTO statuses (path, cpu, memory, created_at)
	VALUES (?, ?, ?, ?);
	`
	_, err := db.d.Exec(insertQuery, path, cpu, memory, time.Now())
	if err != nil {
		log.Println(err)
	}
}
func (db DB) AddRequest(req *http.Request) {
	insertQuery := `
	INSERT INTO requests (ip, url, created_at, is_websocket)
	VALUES (?, ?, ?, ?);
	`
	_, err := db.d.Exec(insertQuery, getIP(req),
		req.URL.Path, time.Now(), isWebSocketRequest(req))
	if err != nil {
		log.Println(err)
	}
}
