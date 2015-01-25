package main

import (
	"log"

	"github.com/gophergala/go_ne/core"
	"strconv"
)

// BUG(Tobscher) use command line arguments to perform correct task
func main() {
	log.SetPrefix("[go-ne] ")

	config, err := core.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = config.Load("config/test-tasks.yaml")
	if err != nil {
		log.Fatal(err)
	}
	
	w, err := NewWeb(); if err != nil {
		log.Fatal(err)
	}

	portString := config.Interfaces["web"]["port"]
	port, err := strconv.ParseUint(portString, 10, 0)
	if(portString == "" || err != nil) {
		port = 20000	// Default
	}
	
	log.Printf("Web interface serving on port %d", port)	
	w.Serve(uint(port), config)
}
