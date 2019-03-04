package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/sector-f/eggchan"
)

type HomePage struct {
	template *template.Template
	modified time.Time
}

func (p HomePage) Route() string {
	return "/"
}

func (p HomePage) Template() string {
	return "index.html"
}

func (p *HomePage) SetTemplate(t *template.Template) {
	p.template = t
}

func (p *HomePage) SetTime(t time.Time) {
	p.modified = t
}

func (p HomePage) GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := http.Get(apiEndpoint + "/boards")
		if err != nil {
			return // TODO: add actual error handling
		}
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		boards := []eggchan.Board{}
		json.Unmarshal(body, &boards)

		fileInfo, err := os.Stat(p.Template())
		if err == nil {
			if fileInfo.ModTime().After(p.modified) {
				newTemplate, err := template.ParseFiles(p.Template())
				if err == nil {
					p.SetTime(fileInfo.ModTime())
					p.SetTemplate(newTemplate)
				}
			}
		}

		p.template.Execute(w, boards)
	}
}

func (p HomePage) PostHandler() http.HandlerFunc {
	return nil
}
