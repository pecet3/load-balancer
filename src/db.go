package main

import (
	"database/sql"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	d *sql.DB
}

func NewDB() (*DB, error) {
	err := os.MkdirAll("database", os.ModePerm)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", "file:database/store.db?_journal_mode=WAL")
	if err != nil {
		return nil, err
	}

	createRequestsTable := `
	CREATE TABLE IF NOT EXISTS requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT NOT NULL,
		body TEXT NOT NULL,
		url TEXT NOT NULL,
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
func (db DB) AddStatus(url string, cpu float64, memory float64) {
	insertQuery := `
	INSERT INTO statuses (url, cpu, memory, created_at)
	VALUES (?, ?, ?, ?);
	`
	_, err := db.d.Exec(insertQuery, url, cpu, memory, time.Now())
	if err != nil {
		log.Println(err)
	}
}
func (db DB) AddRequest(req *http.Request) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
	}
	defer req.Body.Close()

	bodyStr := string(bodyBytes)

	insertQuery := `
	INSERT INTO requests (ip, body, url, created_at)
	VALUES (?, ?, ?, ?);
	`
	_, err = db.d.Exec(insertQuery, getIP(req), bodyStr, req.URL.Path, time.Now())
	if err != nil {
		log.Println(err)
	}
}

func getIP(r *http.Request) string {
	xri := r.Header.Get("X-Real-Ip")
	if xri != "" {
		ip := xri
		if idx := net.ParseIP(ip); idx != nil {
			return ip
		}
	}
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ip := xff
		if idx := net.ParseIP(ip); idx != nil {
			return ip
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		return ""
	}
	return userIP.String()
}
