package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sector-f/eggchan"
)

type ThreadPage struct {
	template *template.Template
}

func (p ThreadPage) Route() string {
	return "/{board}/{thread}"
}

func (p ThreadPage) Template() string {
	return "thread.html"
}

func (p *ThreadPage) SetTemplate(t *template.Template) {
	p.template = t
}

func (p ThreadPage) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		boardName, exists := vars["board"]
		if !exists {
			return
		}

		threadID, exists := vars["thread"]
		if !exists {
			return
		}

		response, err := http.Get(apiEndpoint + "/boards/" + boardName + "/" + threadID)
		if err != nil {
			return // TODO: add actual error handling
		}
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		board := eggchan.ThreadReply{}
		json.Unmarshal(body, &board)
		p.template.Execute(w, board)
	}
}
