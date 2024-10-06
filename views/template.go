package views

import (
	"fmt"
	"html/template"
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
		},
	)

	writer.Header().Set("Content-Type", "text/html charset=UTF-8")

	executeErrpr := htmlTemplate.Execute(writer, data)

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
