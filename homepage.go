package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/sector-f/eggchan"
)

type HomePage struct {
	template *template.Template
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

func (p HomePage) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := http.Get(apiEndpoint + "/boards")
		if err != nil {
			return // TODO: add actual error handling
		}
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		boards := []eggchan.Board{}
		json.Unmarshal(body, &boards)
		p.template.Execute(w, boards)
	}
}
