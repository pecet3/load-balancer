package main

import (
	"database/sql"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	d *sql.DB
}

func NewDB() (*DB, error) {
	db, err := sql.Open("sqlite3", "file:database/requests.db?_journal_mode=WAL")
	if err != nil {
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS requests_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT NOT NULL,
		body TEXT NOT NULL,
		url TEXT NOT NULL,
		created_at DATETIME NOT NULL
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return &DB{d: db}, nil
}

func (db DB) AddRequest(req *http.Request) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
	}
	defer req.Body.Close()

	bodyStr := string(bodyBytes)

	insertQuery := `
	INSERT INTO requests_logs (ip, body, url, created_at)
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
