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

const WEB_FOLDER = "/gokiss"
const STATIC_CACHE_SECONDS = 3600
const WS_READ_BUFFER_SIZE = 1024


type Web struct {
	tplSet   *pongo2.TemplateSet 
}


type WsRequest struct {
	Action   string
	Args     map[string]string
}

type WsResponse struct {
	Type   string              `json:"type"`
	Data   map[string]string   `json:"data"`
}


func NewWeb() (*Web, error) {
	w := Web{}
	
	w.tplSet = pongo2.NewSet("www")
	err := w.tplSet.SetBaseDirectory("www-views/"); if err != nil {
		return nil, err
	}
	
	return &w, nil
}


func (web *Web) Serve(port uint, config *core.Config) (error) {
	mux := routes.New()
	
	// Serve static assets...
	http.Handle(WEB_FOLDER + "/static/", maxAgeHandler(
		STATIC_CACHE_SECONDS,
		http.StripPrefix(WEB_FOLDER + "/static/", http.FileServer(http.Dir("./www-static")))))
	
	// Serve web...
	mux.Get(WEB_FOLDER + "/", web.wwwGokiss())
	mux.Get(WEB_FOLDER + "/task", web.wwwGokissTasks(config))
	mux.Get(WEB_FOLDER + "/task/:taskName/runstatic", web.wwwGokissTaskRunStatic(config))
	mux.Get(WEB_FOLDER + "/task/:taskName/run", web.wwwGokissTaskRun(config))
	mux.Get(WEB_FOLDER + "/auth/log-in", web.wwwGokissAuthLogin())
    http.Handle(WEB_FOLDER + "/socket", websocket.Handler(web.sockGokissTaskRun(config)))

	http.Handle("/", mux)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	
	return nil
}


func (web *Web) wwwGokiss() func(w http.ResponseWriter, r *http.Request) {
	tpl, err := web.tplSet.FromCache("pages/index.pongo"); if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		out, err := tpl.Execute(pongo2.Context{
			"webfolder": WEB_FOLDER,
			"title": "Overview",
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
			"webfolder": WEB_FOLDER,
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
			"webfolder": WEB_FOLDER,
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
			"webfolder": WEB_FOLDER,
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
	task, ok := config.Tasks[taskName]; if !ok {
		io.WriteString(w, "Error!")
		return
	}
	
	io.WriteString(w, fmt.Sprintf("Task %s has %d steps:\n", taskName, len(task.Steps)))
	
	stdOut := make(chan []byte)
	stdErr := make(chan []byte)
	go func() {
		for {
			select {
				case out := <-stdOut:
					if(len(out) == 0) {
//						break;
					}
					
					outString := fmt.Sprintf("%s", out)
					log.Print("OUT: " + outString)

					// #TODO: Handle error...
					web.sendResponseToSocket(w, WsResponse{
						Type: "out",
						Data: map[string]string{
							"message": outString,
						},
					})

				case out := <-stdErr:
					if(len(out) == 0) {
//						break;
					}
				
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

	for _, step := range task.Steps {
		task, _ := core.NewTask(step.Command, step.Args)
		runner, _ := NewLocalRunner()

		err := runner.Run(task, stdOut, stdErr); if err != nil {
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