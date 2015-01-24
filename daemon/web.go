package main

import (
	"net/http"
	"fmt"
	"github.com/drone/routes"
	"github.com/gophergala/go_ne/core"
	"github.com/flosch/pongo2"
	"log"
	"io"
)

const WEB_FOLDER = "/gokiss"
const STATIC_CACHE_SECONDS = 3600


type Web struct {
	tplSet   *pongo2.TemplateSet 
}


func NewWeb() (*Web, error) {
	w := Web{}
	
	w.tplSet = pongo2.NewSet("www")
	err := w.tplSet.SetBaseDirectory("www-views/"); if err != nil {
		return nil, err
	}
	
	return &w, nil
}


func (web* Web) Serve(port uint, config *core.Config) (error) {
	mux := routes.New()

	// Serve static assets...
	http.Handle("/static/", maxAgeHandler(
		STATIC_CACHE_SECONDS,
		http.StripPrefix("/static/", http.FileServer(http.Dir("./www-static")))))

	// Serve web...
	mux.Get(WEB_FOLDER + "/", web.wwwGokiss())
	mux.Get(WEB_FOLDER + "/tasks", web.wwwGokissTasks(config))
	mux.Get(WEB_FOLDER + "/auth/log-in", web.wwwGokissAuthLogin())

	http.Handle("/", mux)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	
	return nil
}


func (web* Web) wwwGokiss() func(w http.ResponseWriter, r *http.Request) {
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


func (web* Web) wwwGokissTasks(config *core.Config) func(w http.ResponseWriter, r *http.Request) {
	tpl, err := web.tplSet.FromCache("pages/tasks.pongo"); if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		out, err := tpl.Execute(pongo2.Context{
			"webfolder": WEB_FOLDER,
			"title": "Tasks",
		}); if err != nil {
			web.errorHandler(w, http.StatusInternalServerError)
			return
		}
		
		io.WriteString(w, out)		
	}
}



func (web* Web) wwwGokissAuthLogin() func(w http.ResponseWriter, r *http.Request) {
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



// For static content...
func maxAgeHandler(seconds int, h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", seconds))
		h.ServeHTTP(w, r)
	})
}


func (web* Web) errorHandler(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}