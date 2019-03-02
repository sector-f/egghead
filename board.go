package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sector-f/eggchan"
)

type BoardPage struct {
	template *template.Template
}

func (p BoardPage) Route() string {
	return "/{board}"
}

func (p BoardPage) Template() string {
	return "board.html"
}

func (p *BoardPage) SetTemplate(t *template.Template) {
	p.template = t
}

func (p BoardPage) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		boardName, exists := vars["board"]
		if !exists {
			return
		}
		response, err := http.Get(apiEndpoint + "/boards/" + boardName)
		if err != nil {
			return // TODO: add actual error handling
		}
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		board := eggchan.BoardReply{}
		json.Unmarshal(body, &board)
		p.template.Execute(w, board)
	}
}
