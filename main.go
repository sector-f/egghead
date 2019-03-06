package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var bindAddr, apiEndpoint, webroot string

type Page interface {
	Route() string
	Template() string
	SetTemplate(*template.Template)
	SetTime(time.Time)
	GetHandler() http.HandlerFunc
	PostHandler() http.HandlerFunc
}

func main() {
	flag.StringVar(&bindAddr, "bind", "127.0.0.1:8080", "Address to bind to")
	flag.StringVar(&apiEndpoint, "api", "http://127.0.0.1:8000", "Backend address to connect to")
	flag.StringVar(&webroot, "webroot", "./", "Directory to look for web files in")
	flag.Parse()

	pages := []Page{
		&HomePage{},
		&BoardPage{},
		&ThreadPage{},
	}

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", handlers.LoggingHandler(os.Stdout, http.FileServer(http.Dir(webroot+"static")))))

	for _, page := range pages {
		fmt.Printf("Loading %s...", page.Template())

		template, err := template.ParseFiles(webroot + page.Template())
		if err != nil {
			fmt.Printf("failed: %s\n", err)
			continue
		}

		fileInfo, err := os.Stat(page.Template())
		if err != nil {
			fmt.Printf("failed: could not stat file\n")
			continue
		}
		page.SetTime(fileInfo.ModTime())

		fmt.Printf("succeeded\n")
		page.SetTemplate(template)

		if page.GetHandler() != nil {
			r.Methods("GET").Path(page.Route()).Handler(handlers.LoggingHandler(os.Stdout, page.GetHandler()))
		}

		if page.PostHandler() != nil {
			r.Methods("POST").Path(page.Route()).Handler(handlers.LoggingHandler(os.Stdout, page.PostHandler()))
		}
	}

	fmt.Println("Now listening on", bindAddr)
	http.ListenAndServe(bindAddr, r)
}
