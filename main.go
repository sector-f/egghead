package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	// "github.com/sector-f/eggchan"
)

var apiEndpoint string = "http://127.0.0.1:8000"

type Page interface {
	Route() string
	Template() string
	SetTemplate(*template.Template)
	Handler() http.HandlerFunc
}

func main() {
	pages := []Page{
		&HomePage{},
		&BoardPage{},
	}

	r := mux.NewRouter()

	for _, page := range pages {
		template, err := template.ParseFiles(page.Template())
		if err != nil {
			fmt.Println("Failed to load ", page.Template(), ": ", err)
			continue
		}
		page.SetTemplate(template)

		r.HandleFunc(
			page.Route(),
			page.Handler(),
		)
	}

	http.ListenAndServe(":8080", r)
}
