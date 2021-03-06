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
		response, err := http.Get(apiEndpoint + "/categories")
		if err != nil {
			w.WriteHeader(http.StatusGatewayTimeout)
			errMessage := errorMessage{http.StatusGatewayTimeout, "Gateway error"}
			errorTemplate.Execute(w, errMessage)
			return
		}
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)

		switch response.StatusCode {
		case http.StatusOK:
			break
		default:
			reply := errorReply{}
			json.Unmarshal(body, &reply)

			w.WriteHeader(response.StatusCode)
			errMessage := errorMessage{response.StatusCode, reply.Error}
			errorTemplate.Execute(w, errMessage)
			return
		}

		categories := []eggchan.Category{}
		json.Unmarshal(body, &categories)

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

		p.template.Execute(w, categories)
	}
}

func (p HomePage) PostHandler() http.HandlerFunc {
	return nil
}
