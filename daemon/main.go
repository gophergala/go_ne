package main

import (
	"log"

	"github.com/gophergala/go_ne/core"
	"strconv"
	"errors"
)

// BUG(Tobscher) use command line arguments to perform correct task
func main() {
	log.SetPrefix("[go-ne] ")

	config, err := core.NewConfig(); if err != nil {
		log.Fatal(err)
	}

	err = config.Load("config/test-tasks.yaml"); if err != nil {
		log.Fatal(err)
	}
	
	webFolder := config.Interfaces.Web.Settings["folder"]
	if(webFolder == "") {
		webFolder = "/gokiss"	// Default
	}

	sessionSecret := config.Interfaces.Web.Settings["sessionsecret"]
	if(sessionSecret == "") {
		log.Fatal(errors.New("Please set a sessionsecret for the web interface"))
	}

	w, err := NewWeb(webFolder, sessionSecret); if err != nil {
		log.Fatal(err)
	}

	_, err = NewEventsMonitor(config, w); if err != nil {
		log.Fatal(err)
	}
	
	portString := config.Interfaces.Web.Settings["port"]
	port, err := strconv.ParseUint(portString, 10, 0)
	if(portString == "" || err != nil) {
		port = 20000	// Default
	}
	
	log.Printf("Web interface serving on port %d", port)	
	w.Serve(uint(port), config)
}
