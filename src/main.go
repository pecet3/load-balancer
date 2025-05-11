package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	log.Println("Starting a loadbalancer :8080...")
	cfg, err := GetConfig()
	if err != nil {
		log.Fatalln("Error with loading a config file!")
	}
	log.Println(cfg.Port)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Błąd serwera: ", err)
	}

}
