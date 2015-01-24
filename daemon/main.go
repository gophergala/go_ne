package main

import (
	"github.com/gophergala/go_ne/core"
	"log"
)

func main() {
	config, err := core.NewConfig(); if err != nil {
		log.Fatal(err)
	}

	err = config.Load("config/test-tasks.yaml"); if err != nil {
		log.Fatal(err)
	}
	
	w, err := NewWeb(); if err != nil {
		log.Fatal(err)
	}

	w.Serve(20000, config)
}
