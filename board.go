package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sector-f/eggchan"
)

type BoardPage struct {
	template *template.Template
	modified time.Time
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

func (p *BoardPage) SetTime(t time.Time) {
	p.modified = t
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

		p.template.Execute(w, board)
	}
}
