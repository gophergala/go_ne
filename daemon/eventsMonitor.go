package main

import(
	"github.com/gophergala/go_ne/core"
	"log"
	"time"
	"fmt"
	"net/http"
	"io"
	"encoding/json"
)


type EventsMonitor struct {
}


func NewEventsMonitor(config *core.Config, web *Web) (*EventsMonitor, error) {
	em := EventsMonitor{}
	
	for triggerName, event := range config.Triggers {
		log.Printf("Loading trigger [%s]\n", triggerName)
		
		switch event.Type {
			case `webhook`:
				err := em.createWebHook(config, triggerName, event, web); if err != nil {
					return nil, err
				}

			case `github-webhook`:
				err := em.createGitHubWebHook(config, triggerName, event, web); if err != nil {
					return nil, err
				}
				
			case `periodic`:
				err := em.createPeriodic(config, triggerName, event); if err != nil {
					return nil, err
				}
		}
	}
	
	return &em, nil
}



func (em *EventsMonitor) createWebHook(config *core.Config, triggerName string, event core.ConfigEvent, web *Web) (error) {
	web.mux.Get(web.webFolder + "/triggers/" + event.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		err := em.runTask(config, event.Task); if err != nil {
			log.Printf("Webhook [start]: %s for %s\n", triggerName, event.ServerGroup)

			// #TODO: Persistent Logging
			log.Printf("%s", err)
			
			log.Printf("WebHook [failed]: %s for %s\n", triggerName, event.ServerGroup)
			
			// #TODO: Server error response?
			return
		}

		log.Printf("Webhook [complete]: %s for %s\n", triggerName, event.ServerGroup)
		io.WriteString(w, "OK")
	})

	return nil
}


type githubHookData struct {
	Ref string
}


//'push'
func (em *EventsMonitor) createGitHubWebHook(config *core.Config, triggerName string, event core.ConfigEvent, web *Web) (error) {
	web.mux.Post(web.webFolder + "/triggers/" + event.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%+v", r.Header)

		if(r.Header.Get("Content-Type") != `application/json` || r.Header.Get("X-Github-Event") != `push`) {
			return
		}
		
		log.Printf("\n\n%+v\n", r.Body)

		decoder := json.NewDecoder(r.Body)

		var t githubHookData
		err := decoder.Decode(&t); if err != nil {
			fmt.Print(err)
		}

		log.Printf("%+v", t)

		err = em.runTask(config, event.Task); if err != nil {
			log.Printf("GitHub Webhook [start]: %s for %s\n", triggerName, event.ServerGroup)

			// #TODO: Persistent Logging
			log.Printf("%s", err)
			
			log.Printf("GitHub WebHook [failed]: %s for %s\n", triggerName, event.ServerGroup)
			
			// #TODO: Server error response?
			return
		}

		log.Printf("GitHub Webhook [complete]: %s for %s\n", triggerName, event.ServerGroup)
		io.WriteString(w, "OK")
	})

	return nil
}



func (em *EventsMonitor) createPeriodic(config *core.Config, triggerName string, event core.ConfigEvent) (error) {
	ticker := time.NewTicker(time.Duration(event.Period) * time.Second)
		
	go func() {
		for {
			select {
				case <- ticker.C:
					log.Printf("Periodic [start]: %s for %s\n", triggerName, event.ServerGroup)
					
					err := em.runTask(config, event.Task); if err != nil {
						// #TODO: Persistent Logging
						log.Printf("%s", err)
						
						log.Printf("Periodic [failed]: %s for %s\n", triggerName, event.ServerGroup)
						continue
					}
					
					log.Printf("Periodic [complete]: %s for %s\n", triggerName, event.ServerGroup)
			}		
		}
	}()

	return nil
}


func (em *EventsMonitor) runTask(config *core.Config, task string) (error) {
	runner, err := NewLocalRunner(); if err != nil {
		return err
	}

	// #TODO: improve [/dev/null it for the moment!]
	go func() {
		for {
			select {
				case out := <-runner.chStdOut:
					log.Printf("%s", out)
				case out := <-runner.chStdErr:
					log.Printf("%s", out)
			}
		}
	}()

	err = core.RunTask(runner, config, task); if err != nil {
		return err
	}
	
	return nil
}
