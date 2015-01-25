package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/drone/routes"
	"github.com/flosch/pongo2"
	"github.com/gophergala/go_ne/core"
	"github.com/gorilla/sessions"
	"golang.org/x/net/websocket"
)

const STATIC_CACHE_SECONDS = 3600
const WS_READ_BUFFER_SIZE = 1024

type Web struct {
	tplSet       *pongo2.TemplateSet
	mux          *routes.RouteMux
	webFolder    string
	sessionStore *sessions.CookieStore
}


type WsRequest struct {
	Action string
	Args   map[string]string
}

type WsResponse struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

func NewWeb(webFolder string, sessionSecret string) (*Web, error) {
	web := Web{
		webFolder: webFolder,
		sessionStore: sessions.NewCookieStore([]byte(sessionSecret)),
	}

	web.tplSet = pongo2.NewSet("www")
	err := web.tplSet.SetBaseDirectory("www-views/")
	if err != nil {
		return nil, err
	}

	web.mux = routes.New()

	return &web, nil
}

func (web *Web) Serve(port uint, config *core.Config) error {
	// Serve static assets...
	http.Handle(web.webFolder+"/static/", maxAgeHandler(
		STATIC_CACHE_SECONDS,
		http.StripPrefix(web.webFolder+"/static/", http.FileServer(http.Dir("./www-static")))))

	// Serve web...
	web.mux.Get(web.webFolder+"/about", web.wwwGokissAbout())
	web.mux.Get(web.webFolder+"/auth/log-in", web.wwwGokissAuthLogin(config))
	web.mux.Post(web.webFolder+"/auth/log-in", web.wwwPostGokissAuthLogin(config))
	web.mux.Get(web.webFolder+"/auth/log-out", web.wwwGokissAuthLogout(config))
	web.mux.Get(web.webFolder+"/", web.wwwGokiss(config))
	web.mux.Get(web.webFolder+"/servergroup", web.wwwGokissServergroups(config))
	web.mux.Get(web.webFolder+"/servergroup/:servergroupName", web.wwwGokissServergroup(config))
	web.mux.Get(web.webFolder+"/task", web.wwwGokissTasks(config))
	web.mux.Get(web.webFolder+"/servergroup/:servergroupName/task/:taskName/run", web.wwwGokissTaskRun(config))
	http.Handle(web.webFolder+"/socket", websocket.Handler(web.sockGokissTaskRun(config)))

	http.Handle("/", web.mux)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	return nil
}


func (web *Web) wwwGokissAbout() func(w http.ResponseWriter, r *http.Request) {
	wh, err := web.newWebHandler(); if err != nil {
		log.Fatal(err)
	}

	wh.setCheckAuth(false)
	wh.setupTemplate(web, "pages/about.pongo")

	return func(w http.ResponseWriter, r *http.Request) {
		wh.newRequest()
		s := wh.getSession(web, r)
		defer wh.output(web, w, r, s)
				
		if(wh.canShowPage(web, s)) {
			wh.renderTemplate(pongo2.Context{
				"webfolder": web.webFolder,
				"title":     "About",
				"auth": wh.hasAuth(web, s),
			})
		}
	}
}




func (web *Web) wwwGokissAuthLogin(config *core.Config) func(w http.ResponseWriter, r *http.Request) {
	wh, err := web.newWebHandler(); if err != nil {
		log.Fatal(err)
	}

	wh.setCheckAuth(false)
	wh.setupTemplate(web, "pages/auth/log-in.pongo")

	return func(w http.ResponseWriter, r *http.Request) {
		wh.newRequest()
		s := wh.getSession(web, r)
		defer wh.output(web, w, r, s)
				
		if(wh.canShowPage(web, s)) {
			wh.renderTemplate(pongo2.Context{
				"webfolder": web.webFolder,
				"title":     "Log In",
				"auth": wh.hasAuth(web, s),
			})
		}
	}
}


func (web *Web) wwwPostGokissAuthLogin(config *core.Config) func(w http.ResponseWriter, r *http.Request) {
	wh, err := web.newWebHandler(); if err != nil {
		log.Fatal(err)
	}

	wh.setCheckAuth(false)

	return func(w http.ResponseWriter, r *http.Request) {
		wh.newRequest()
		s := wh.getSession(web, r)
		defer wh.output(web, w, r, s)
		
		username := r.FormValue("username")
		password := r.FormValue("password")

		authorised := false
		for _, user := range config.Interfaces.Web.Users {
			if(username == user.Username) {
				if(password == user.Password) {
					authorised = true
				}
				
				break
			}
		}
		
		if(!authorised) {
			wh.setRedirect(web.webFolder + "/auth/log-in")
			return
		}
		
		s.Values[`auth`] = true
		wh.setRedirect(web.webFolder + "/")
	}
}



func (web *Web) wwwGokissAuthLogout(config *core.Config) func(w http.ResponseWriter, r *http.Request) {
	wh, err := web.newWebHandler(); if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		wh.newRequest()
		s := wh.getSession(web, r)
		defer wh.output(web, w, r, s)

		delete(s.Values, `auth`)
		wh.setRedirect(web.webFolder + "/auth/log-in")
	}
}



func (web *Web) wwwGokiss(config *core.Config) func(w http.ResponseWriter, r *http.Request) {
	wh, err := web.newWebHandler(); if err != nil {
		log.Fatal(err)
	}

	wh.setCheckAuth(true)
	wh.setupTemplate(web, "pages/index.pongo")

	return func(w http.ResponseWriter, r *http.Request) {
		wh.newRequest()
		s := wh.getSession(web, r)
		defer wh.output(web, w, r, s)
				
		if(wh.canShowPage(web, s)) {
			wh.renderTemplate(pongo2.Context{
				"webfolder":    web.webFolder,
				"title":        "Overview",
				"servergroups": config.ServerGroups,
				"auth": wh.hasAuth(web, s),
			})
		}
	}
}



func (web *Web) wwwGokissServergroups(config *core.Config) func(w http.ResponseWriter, r *http.Request) {
	wh, err := web.newWebHandler(); if err != nil {
		log.Fatal(err)
	}

	wh.setCheckAuth(true)
	wh.setupTemplate(web, "pages/servergroups.pongo")

	return func(w http.ResponseWriter, r *http.Request) {
		wh.newRequest()
		s := wh.getSession(web, r)
		defer wh.output(web, w, r, s)
				
		if(wh.canShowPage(web, s)) {
			wh.renderTemplate(pongo2.Context{
				"webfolder": web.webFolder,
				"title":        "Overview",
				"servergroups": config.ServerGroups,
				"auth": wh.hasAuth(web, s),
			})
		}
	}
}


func (web *Web) wwwGokissServergroup(config *core.Config) func(w http.ResponseWriter, r *http.Request) {
	wh, err := web.newWebHandler(); if err != nil {
		log.Fatal(err)
	}

	wh.setCheckAuth(true)
	wh.setupTemplate(web, "pages/servergroup.pongo")

	return func(w http.ResponseWriter, r *http.Request) {
		wh.newRequest()
		s := wh.getSession(web, r)
		defer wh.output(web, w, r, s)
				
		if(wh.canShowPage(web, s)) {
			params := r.URL.Query()
			servergroupName := params.Get(":servergroupName")

			servergroup, ok := config.ServerGroups[servergroupName]
			if !ok {
				web.errorHandler(w, http.StatusNotFound)
				return
			}

			wh.renderTemplate(pongo2.Context{
				"webfolder":       web.webFolder,
				"title":           "Overview",
				"servergroupName": servergroupName,
				"servergroup":     servergroup,
				"tasks":           config.Tasks,
				"auth": wh.hasAuth(web, s),
			})
		}
	}
}


func (web *Web) wwwGokissTasks(config *core.Config) func(w http.ResponseWriter, r *http.Request) {
	wh, err := web.newWebHandler(); if err != nil {
		log.Fatal(err)
	}

	wh.setCheckAuth(true)
	wh.setupTemplate(web, "pages/tasks.pongo")

	return func(w http.ResponseWriter, r *http.Request) {
		wh.newRequest()
		s := wh.getSession(web, r)
		defer wh.output(web, w, r, s)
				
		if(wh.canShowPage(web, s)) {
			wh.renderTemplate(pongo2.Context{
				"webfolder": web.webFolder,
				"title":     "Tasks",
				"tasks":     config.Tasks,
				"auth": wh.hasAuth(web, s),
			})
		}
	}
}



func (web *Web) wwwGokissTaskRun(config *core.Config) func(w http.ResponseWriter, r *http.Request) {
	wh, err := web.newWebHandler(); if err != nil {
		log.Fatal(err)
	}

	wh.setCheckAuth(true)
	wh.setupTemplate(web, "pages/task-run.pongo")

	return func(w http.ResponseWriter, r *http.Request) {
		wh.newRequest()
		s := wh.getSession(web, r)
		defer wh.output(web, w, r, s)
				
		params := r.URL.Query()
		servergroupName := params.Get(":servergroupName")
		taskName := params.Get(":taskName")

		_, ok := config.ServerGroups[servergroupName]
		if !ok {
			web.errorHandler(w, http.StatusNotFound)
			return
		}

		_, ok = config.Tasks[taskName]
		if !ok {
			web.errorHandler(w, http.StatusNotFound)
			return
		}

		if(wh.canShowPage(web, s)) {
			wh.renderTemplate(pongo2.Context{
				"host":      r.Host,
				"webfolder": web.webFolder,
				"title":     "Run Task",
				"servergroupName":  servergroupName,
				"taskName":  taskName,
				"auth": wh.hasAuth(web, s),
			})
		}
	}
}



// # Websocket...

func (web *Web) sockGokissTaskRun(config *core.Config) func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {
		wh, err := web.newWebHandler(); if err != nil {
			log.Fatal(err)
		}

		wh.newRequest()
		s := wh.getSession(web, ws.Request())

		if(!wh.hasAuth(web, s)) {
			// #TODO: Respond with appropriate HTTP code
			return
		}

		for {
			message := make([]byte, WS_READ_BUFFER_SIZE)
			length, err := ws.Read(message)
			if err != nil {
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
					// request.Args["servergroupName"]
					web.taskRun(ws, request.Args["taskName"], config)
			}
		}
	}
}


func (web *Web) taskRun(w io.Writer, taskName string, config *core.Config) {
	runner, err := core.NewLocalRunner()
	if err != nil {
		io.WriteString(w, "Error!")
		return
	}

	go func() {
		for {
			select {
			case out := <-runner.ChStdOut:
				outString := fmt.Sprintf("%s", out)
				log.Print("OUT: " + outString)

				// #TODO: Handle error...
				web.sendResponseToSocket(w, WsResponse{
					Type: "out",
					Data: map[string]string{
						"message": outString,
					},
				})

			case out := <-runner.ChStdErr:
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

	err = core.RunTask(runner, config, taskName)
	if err != nil {
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

func (web *Web) sendResponseToSocket(w io.Writer, r WsResponse) error {
	b, err := json.Marshal(r)
	if err != nil {
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



type WebHandler struct {
	checkAuth     bool
	template      *pongo2.Template
	body          string
	redirectUrl   string
}


func (web *Web) newWebHandler() (*WebHandler, error) {
	wh := WebHandler{
		checkAuth: false,
	}
	
	return &wh, nil
}


// #TODO - this is a sledgehammer for a bug - the webhandler should ideally be refactored into persistent handler and request handler
func (wh *WebHandler) newRequest() {
	wh.body = ""
	wh.redirectUrl = ""
}


func (wh *WebHandler) setCheckAuth(checkAuth bool) {
	wh.checkAuth = checkAuth
}


func (wh *WebHandler) setRedirect(redirectUrl string) {
	wh.redirectUrl = redirectUrl
}


func (wh *WebHandler) setupTemplate(web *Web, template string) {
	var err error
	wh.template, err = web.tplSet.FromCache(template); if err != nil {
		log.Fatal(err)
	}
}


func (wh *WebHandler) getSession(web *Web, r *http.Request) *sessions.Session {
	session, _ := web.sessionStore.Get(r, "gokiss")
	return session
}


func (wh *WebHandler) hasAuth(web *Web, s *sessions.Session) bool {
	auth, ok := s.Values[`auth`]; if !ok || auth != true {
		return false
	}
	
	return true
}


func (wh *WebHandler) canShowPage(web *Web, s *sessions.Session) bool {
	if(wh.checkAuth && !wh.hasAuth(web, s)) {
		wh.setRedirect(web.webFolder + "/auth/log-in")
		return false
	}
	
	return true
}


func (wh *WebHandler) renderTemplate(context pongo2.Context) {
	var err error
	wh.body, err = wh.template.Execute(context); if err != nil {
		//wh.errorHandler(w, http.StatusInternalServerError)
	}
}


func (wh *WebHandler) output(web *Web, w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	s.Save(r, w)
	
	if(wh.redirectUrl != "") {
		http.Redirect(w, r, wh.redirectUrl, http.StatusFound)
	} else {	
		io.WriteString(w, wh.body)
	}
}
