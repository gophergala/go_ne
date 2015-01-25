package main

import (
	"net/http"
	"fmt"
	"github.com/drone/routes"
	"github.com/gophergala/go_ne/core"
	"github.com/flosch/pongo2"
	"log"
	"io"
	"golang.org/x/net/websocket"
	"encoding/json"
)

const STATIC_CACHE_SECONDS = 3600
const WS_READ_BUFFER_SIZE = 1024


type Web struct {
	tplSet      *pongo2.TemplateSet
	mux         *routes.RouteMux
	webFolder   string
}


type WsRequest struct {
	Action   string
	Args     map[string]string
}

type WsResponse struct {
	Type   string              `json:"type"`
	Data   map[string]string   `json:"data"`
}


func NewWeb(webFolder string) (*Web, error) {
	w := Web{
		webFolder: webFolder,
	}
	
	w.tplSet = pongo2.NewSet("www")
	err := w.tplSet.SetBaseDirectory("www-views/"); if err != nil {
		return nil, err
	}
	
	w.mux = routes.New()
	
	return &w, nil
}


func (web *Web) Serve(port uint, config *core.Config) (error) {
	// Serve static assets...
	http.Handle(web.webFolder + "/static/", maxAgeHandler(
		STATIC_CACHE_SECONDS,
		http.StripPrefix(web.webFolder + "/static/", http.FileServer(http.Dir("./www-static")))))
	
	// Serve web...
	web.mux.Get(web.webFolder + "/", web.wwwGokiss(config))
	web.mux.Get(web.webFolder + "/about", web.wwwGokissAbout())
	web.mux.Get(web.webFolder + "/servergroup", web.wwwGokissServergroups(config))
	web.mux.Get(web.webFolder + "/servergroup/:servergroupName", web.wwwGokissServergroup(config))
	web.mux.Get(web.webFolder + "/task", web.wwwGokissTasks(config))
	web.mux.Get(web.webFolder + "/servergroup/:servergroupName/task/:taskName/runstatic", web.wwwGokissTaskRunStatic(config))
	web.mux.Get(web.webFolder + "/servergroup/:servergroupName/task/:taskName/run", web.wwwGokissTaskRun(config))
	web.mux.Get(web.webFolder + "/auth/log-in", web.wwwGokissAuthLogin())
    http.Handle(web.webFolder + "/socket", websocket.Handler(web.sockGokissTaskRun(config)))

	http.Handle("/", web.mux)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	
	return nil
}


func (web *Web) wwwGokiss(config *core.Config) func(w http.ResponseWriter, r *http.Request) {
	tpl, err := web.tplSet.FromCache("pages/index.pongo"); if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		out, err := tpl.Execute(pongo2.Context{
			"webfolder": web.webFolder,
			"title": "Overview",
			"servergroups": config.ServerGroups,
		}); if err != nil {
			web.errorHandler(w, http.StatusInternalServerError)
			return
		}
		
		io.WriteString(w, out)
	}
}



func (web *Web) wwwGokissServergroups(config *core.Config) func(w http.ResponseWriter, r *http.Request) {
	tpl, err := web.tplSet.FromCache("pages/servergroups.pongo"); if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		out, err := tpl.Execute(pongo2.Context{
			"webfolder": web.webFolder,
			"title": "Overview",
			"servergroups": config.ServerGroups,
		}); if err != nil {
			web.errorHandler(w, http.StatusInternalServerError)
			return
		}
		
		io.WriteString(w, out)		
	}
}


func (web *Web) wwwGokissAbout() func(w http.ResponseWriter, r *http.Request) {
	tpl, err := web.tplSet.FromCache("pages/about.pongo"); if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		out, err := tpl.Execute(pongo2.Context{
			"webfolder": web.webFolder,
			"title": "About",
		}); if err != nil {
			web.errorHandler(w, http.StatusInternalServerError)
			return
		}
		
		io.WriteString(w, out)		
	}
}


func (web *Web) wwwGokissServergroup(config *core.Config) func(w http.ResponseWriter, r *http.Request) {
	tpl, err := web.tplSet.FromCache("pages/servergroup.pongo"); if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		servergroupName := params.Get(":servergroupName")

		servergroup, ok := config.ServerGroups[servergroupName]; if !ok {
			web.errorHandler(w, http.StatusNotFound)		
			return
		}
				
		out, err := tpl.Execute(pongo2.Context{
			"webfolder": web.webFolder,
			"title": "Overview",
			"servergroupName": servergroupName,
			"servergroup": servergroup,
			"tasks": config.Tasks,
		}); if err != nil {
			web.errorHandler(w, http.StatusInternalServerError)
			return
		}
		
		io.WriteString(w, out)		
	}
}


func (web *Web) wwwGokissTasks(config *core.Config) func(w http.ResponseWriter, r *http.Request) {
	tpl, err := web.tplSet.FromCache("pages/tasks.pongo"); if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		out, err := tpl.Execute(pongo2.Context{
			"webfolder": web.webFolder,
			"title": "Tasks",
			"tasks": config.Tasks,
		}); if err != nil {
			web.errorHandler(w, http.StatusInternalServerError)
			return
		}
		
		io.WriteString(w, out)		
	}
}



func (web *Web) wwwGokissAuthLogin() func(w http.ResponseWriter, r *http.Request) {
	tpl, err := web.tplSet.FromCache("pages/auth/log-in.pongo"); if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		out, err := tpl.Execute(pongo2.Context{
			"webfolder": web.webFolder,
			"title": "Log In",
		}); if err != nil {
			web.errorHandler(w, http.StatusInternalServerError)
			return
		}
		
		io.WriteString(w, out)
	}
}



func (web *Web) wwwGokissTaskRunStatic(config *core.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		taskName := params.Get(":taskName")

		web.taskRun(w, taskName, config)
	}
}



func (web *Web) wwwGokissTaskRun(config *core.Config) func(w http.ResponseWriter, r *http.Request) {
	tpl, err := web.tplSet.FromCache("pages/task-run.pongo"); if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		taskName := params.Get(":taskName")

		out, err := tpl.Execute(pongo2.Context{
			"host": r.Host,
			"webfolder": web.webFolder,
			"title": "Run Task",
			"taskName": taskName,
		}); if err != nil {
			web.errorHandler(w, http.StatusInternalServerError)
			return
		}
		
		io.WriteString(w, out)
	}
}



// # Websocket...

func (web *Web) sockGokissTaskRun(config *core.Config) func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {
		for {
			message := make([]byte, WS_READ_BUFFER_SIZE)
			length, err := ws.Read(message); if err != nil {
				break
			}
			
			var request WsRequest
			err = json.Unmarshal(message[:length], &request)
			if err != nil {
				/* TODO: Handle gracefully with response to front-end */
				fmt.Println("error:", err)
				continue
			}

			switch request.Action {
				case `task-run`:
					web.taskRun(ws, request.Args["taskName"], config)
			}
		}
	}
}



func (web *Web) taskRun(w io.Writer, taskName string, config *core.Config) {
	runner, err := NewLocalRunner(); if err != nil {
		io.WriteString(w, "Error!")
		return	
	}

	go func() {
		for {
			select {
				case out := <-runner.chStdOut:				
					outString := fmt.Sprintf("%s", out)
					log.Print("OUT: " + outString)

					// #TODO: Handle error...
					web.sendResponseToSocket(w, WsResponse{
						Type: "out",
						Data: map[string]string{
							"message": outString,
						},
					})

				case out := <-runner.chStdErr:
					outString := fmt.Sprintf("%s", out)
					log.Print("ERR: " + outString)

					// #TODO: Handle error...
					web.sendResponseToSocket(w, WsResponse{
						Type: "err",
						Data: map[string]string{
							"message": outString,
						},
					})
			}
		}
	}()
	
	err = core.RunTask(runner, config, taskName); if err != nil {
		outString := fmt.Sprintf("%s", err)
		io.WriteString(w, outString)

		// #TODO: Handle error...
		web.sendResponseToSocket(w, WsResponse{
			Type: "err",
			Data: map[string]string{
				"message": outString,
			},
		})
	}
		
	io.WriteString(w, "Complete!")
}


func (web *Web) sendResponseToSocket(w io.Writer, r WsResponse) (error) {
	b, err := json.Marshal(r); if err != nil {
		return err
	}

	w.Write(b)
	return nil
}


// For static content...
func maxAgeHandler(seconds int, h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", seconds))
		h.ServeHTTP(w, r)
	})
}



func (web *Web) errorHandler(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}