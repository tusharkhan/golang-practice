package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type Template struct {
	htmlTemplate *template.Template
}

func ParseFS(fs fs.FS, pattern ...string) (Template, error) {
	htmlTpl, err := template.ParseFS(fs, pattern...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{
		htmlTemplate: htmlTpl,
	}, nil
}

func (t Template) Execute(writer http.ResponseWriter, data interface{}) {
	writer.Header().Set("Content-Type", "text/html charset=UTF-8")

	executeErrpr := t.htmlTemplate.Execute(writer, data)

	if executeErrpr != nil {
		log.Printf("executing template %v", executeErrpr)
		http.Error(writer, "There was an error executing template", http.StatusInternalServerError)
		return
	}
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
