package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/shirou/gopsutil/v3/process"
)

func main() {
	port := flag.String("port", "8080", "Port for the HTTP server to listen on")
	name := flag.String("name", "MyServer", "Name of the server")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(*name, " is serving html")
		html := fmt.Sprintf(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>%s</title>
				<style>
					body { font-family: Arial, sans-serif; margin: 50px; }
					h1 { color: #333; }
				</style>
			</head>
			<body>
				<h1>Welcome to server: %s</h1>
			</body>
			</html>
		`, *name, *name)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, html)
	})

	http.HandleFunc("/api/statusz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		status, err := GetStatus()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(status)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	addr := fmt.Sprintf(":%s", *port)
	log.Printf("Starting server '%s' on port %s...\n", *name, *port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

type Status struct {
	CPU    float64 `json:"cpu"`
	Memory float32 `json:"memory"`
}

func GetStatus() (*Status, error) {
	pid := os.Getpid()
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil, err
	}
	cpu, err := proc.CPUPercent()
	if err != nil {
		return nil, err
	}
	mem, err := proc.MemoryPercent()
	if err != nil {
		return nil, err
	}
	status := &Status{
		CPU:    cpu,
		Memory: mem,
	}
	return status, nil
}
