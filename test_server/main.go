package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/v3/process"
)

// go run . --port 8083 --name "server 2"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // w prostych testach pozwól na dowolne połączenia
	},
}

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
		type Status struct {
			CPU    float64 `json:"cpu"`
			Memory float32 `json:"memory"`
		}

		pid := os.Getpid()
		proc, err := process.NewProcess(int32(pid))
		if err != nil {
			http.Error(w, "failed to get process", http.StatusInternalServerError)
			return
		}

		cpu, err := proc.CPUPercent()
		if err != nil {
			http.Error(w, "failed to get CPU usage", http.StatusInternalServerError)
			return
		}

		mem, err := proc.MemoryPercent()
		if err != nil {
			http.Error(w, "failed to get memory usage", http.StatusInternalServerError)
			return
		}

		status := Status{
			CPU:    cpu,
			Memory: mem,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(status); err != nil {
			http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}
		defer conn.Close()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			log.Printf("Received: %s", msg)

			response := fmt.Sprintf("Echo from %s: %s", *name, msg)
			if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
				log.Println("Write error:", err)
				break
			}
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
