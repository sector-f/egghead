package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var apiEndpoint string = "http://127.0.0.1:8000"

type Page interface {
	Route() string
	Template() string
	SetTemplate(*template.Template)
	SetTime(time.Time)
	Handler() http.HandlerFunc
}

func main() {
	pages := []Page{
		&HomePage{},
		&BoardPage{},
		&ThreadPage{},
	}

	r := mux.NewRouter()

	for _, page := range pages {
		fmt.Printf("Loading %s...", page.Template())

		template, err := template.ParseFiles(page.Template())
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

		r.HandleFunc(
			page.Route(),
			page.Handler(),
		)
	}

	http.ListenAndServe(":8080", r)
}
