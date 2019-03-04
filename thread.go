package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sector-f/eggchan"
)

type ThreadPage struct {
	template *template.Template
	modified time.Time
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

func (p *ThreadPage) SetTime(t time.Time) {
	p.modified = t
}

func (p ThreadPage) GetHandler() http.HandlerFunc {
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

func (p ThreadPage) PostHandler() http.HandlerFunc {
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

		r.Body = http.MaxBytesReader(w, r.Body, 4096)
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			return
		}

		r.ParseMultipartForm(32 << 20)
		comment := strings.TrimSpace(r.FormValue("comment"))
		if comment == "" {
			return
		}

		author := strings.TrimSpace(r.FormValue("author"))
		if author == "" {
			author = "Anonymous"
		}

		var buf bytes.Buffer
		var client http.Client
		writer := multipart.NewWriter(&buf)

		err = writer.WriteField("author", author)
		if err != nil {
			return
		}

		err = writer.WriteField("comment", comment)
		if err != nil {
			return
		}

		writer.Close()

		request, err := http.NewRequest("POST", apiEndpoint+"/boards/"+boardName+"/"+threadID, &buf)
		if err != nil {
			return
		}
		request.Header.Set("Content-Type", writer.FormDataContentType())

		response, err := client.Do(request)
		if err != nil {
			return
		}

		if response.StatusCode != http.StatusCreated {
			fmt.Println("Error")
			return
		}

		responseMessage, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return
		}

		postResponse := eggchan.PostCommentResponse{}
		err = json.Unmarshal(responseMessage, &postResponse)
		if err != nil {
			return
		}
		http.Redirect(w, r, "/"+boardName+"/"+strconv.Itoa(postResponse.ReplyTo), 303)
	}
}
