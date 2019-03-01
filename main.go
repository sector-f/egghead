package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	_ "github.com/gorilla/mux"
	"github.com/sector-f/eggchan"
)

type page struct {
	template string // Parsed into template.Template
	route    string // Frontend route
	backend  string // API endpoint that gets used
	data     interface{}
}

type boardsData struct {
	BaseUrl string
	Boards  []eggchan.Board
}

func main() {
	apiEndpoint := "http://127.0.0.1:8000"

	pages := []page{
		page{
			template: "index.html",
			route:    "/",
			backend:  "/boards",
			data:     []eggchan.Board{},
		},
	}

	for _, page := range pages {
		template, err := template.ParseFiles(page.template)
		if err != nil {
			fmt.Println("Failed to load ", page.template, ": ", err)
			continue
		}

		http.HandleFunc(
			page.route,
			func(w http.ResponseWriter, r *http.Request) {
				response, err := http.Get(apiEndpoint + page.backend)
				if err != nil {
					return // TODO: add actual error handling
				}
				defer response.Body.Close()

				body, err := ioutil.ReadAll(response.Body)
				json.Unmarshal(body, &page.data)
				template.Execute(w, page.data)
			},
		)
	}

	http.ListenAndServe(":8080", nil)
}
