package main

import (
	"html/template"
)

type errorReply struct {
	Error string
}

type errorMessage struct {
	StatusCode    int
	StatusMessage string
}

var errorTemplate *template.Template

func init() {
	template := template.New("errorTemplate")
	parsed, _ := template.Parse(`<!DOCTYPE html>
<html>
	<head>
		<title>
			{{ .StatusCode }} - {{ .StatusMessage }}
		</title>
		<link rel="stylesheet" href="/static/style.css">
	</head>
	<body>
		<div class="errorcontainer">
			{{ .StatusCode }} - {{ .StatusMessage }}
		</div>
	</body>
</html>`)

	errorTemplate = parsed
}
