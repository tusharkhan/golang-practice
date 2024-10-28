package views

import (
	"bytes"
	"course/context"
	"course/models"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
)

type Template struct {
	htmlTemplate *template.Template
}

func ParseFS(fs fs.FS, pattern ...string) (Template, error) {
	htmlTpl := template.New(pattern[0])

	htmlTpl = htmlTpl.Funcs(
		template.FuncMap{
			"csrf": func() template.HTML {
				return `<! -- tag and ends with a -->`
			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("not implemented")
			},
		},
	)

	htmlTpl, err := htmlTpl.ParseFS(fs, pattern...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{
		htmlTemplate: htmlTpl,
	}, nil
}

func (t Template) Execute(writer http.ResponseWriter, request *http.Request, data interface{}) {
	htmlTemplate, templateCloneError := t.htmlTemplate.Clone()

	if templateCloneError != nil {
		log.Printf("template cloning error %v", templateCloneError)
		http.Error(writer, "Error in rendering template", http.StatusInternalServerError)
		return
	}

	htmlTemplate = htmlTemplate.Funcs(
		template.FuncMap{
			"csrf": func() template.HTML {
				return csrf.TemplateField(request)
			},
			"currentUser": func() *models.User {
				return context.User(request.Context())
			},
		},
	)

	writer.Header().Set("Content-Type", "text/html charset=UTF-8")
	var buffer bytes.Buffer
	executeErrpr := htmlTemplate.Execute(&buffer, data)

	if executeErrpr != nil {
		log.Printf("executing template %v", executeErrpr)
		http.Error(writer, "There was an error executing template", http.StatusInternalServerError)
		return
	}

	io.Copy(writer, &buffer)
}

func ParseTemplate(filepath string) (Template, error) {
	tmp, templateError := template.ParseFiles(filepath)

	if templateError != nil {
		return Template{}, fmt.Errorf("%v", templateError)
	}

	return Template{
		htmlTemplate: tmp,
	}, nil
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}

	return t
}
