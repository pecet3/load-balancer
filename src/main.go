package main

import (
	"log"
)

func main() {
	cfg, err := GetConfig()
	if err != nil {
		log.Fatalln("Error with loading a config file!")
	}
	c := NewCore(cfg)
	c.RunLoadBalancer()
}
