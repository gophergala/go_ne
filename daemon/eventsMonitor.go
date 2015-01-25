package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gophergala/go_ne/core"
)

type EventsMonitor struct {
}

func NewEventsMonitor(config *core.Config, web *Web) (*EventsMonitor, error) {
	em := EventsMonitor{}

	for triggerName, event := range config.Triggers {
		log.Printf("Loading trigger [%s]\n", triggerName)

		switch event.Type {
		case `webhook`:
			err := em.createWebHook(config, triggerName, event, web)
			if err != nil {
				return nil, err
			}

		case `github-webhook`:
			err := em.createGitHubWebHook(config, triggerName, event, web)
			if err != nil {
				return nil, err
			}

		case `periodic`:
			err := em.createPeriodic(config, triggerName, event)
			if err != nil {
				return nil, err
			}
		}
	}

	return &em, nil
}

func (em *EventsMonitor) createWebHook(config *core.Config, triggerName string, event core.ConfigEvent, web *Web) error {
	web.mux.Get(web.webFolder+"/triggers/"+event.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		err := em.runTask(config, event.Task)
		if err != nil {
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

func (em *EventsMonitor) createGitHubWebHook(config *core.Config, triggerName string, event core.ConfigEvent, web *Web) error {
	web.mux.Post(web.webFolder+"/triggers/"+event.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != `application/json` || r.Header.Get("X-Github-Event") != `push` {
			// #TODO: Log
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			// #TODO: Log
			return
		}

		// If we've specified a secret in the YAML then check it...
		if event.Secret != "" {
			requestSecret := r.Header.Get("X-Hub-Signature")
			if requestSecret == "" {
				// #TODO: Log
				return
			}

			// Encode the body with our GitHub secret, and match against the SHA1=secret in the header
			key := []byte(event.Secret)
			h := hmac.New(sha1.New, key)
			h.Write(body)
			secretSha1 := fmt.Sprintf("sha1=%s", hex.EncodeToString(h.Sum(nil)))

			if secretSha1 != requestSecret {
				// #TODO: Log
				return
			}
		}

		var hook githubHookData
		err = json.Unmarshal(body, &hook)
		if err != nil {
			// #TODO: Log
			return
		}

		if hook.Ref != "refs/heads/master" {
			// #TODO: Log
			return
		}

		err = em.runTask(config, event.Task)
		if err != nil {
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

func (em *EventsMonitor) createPeriodic(config *core.Config, triggerName string, event core.ConfigEvent) error {
	ticker := time.NewTicker(time.Duration(event.Period) * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Printf("Periodic [start]: %s for %s\n", triggerName, event.ServerGroup)

				err := em.runTask(config, event.Task)
				if err != nil {
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

func (em *EventsMonitor) runTask(config *core.Config, task string) error {
	runner, err := core.NewLocalRunner()
	if err != nil {
		return err
	}

	// #TODO: improve [/dev/null it for the moment!]
	go core.LogOutput(runner)

	err = core.RunTask(runner, config, task)
	if err != nil {
		return err
	}

	return nil
}
