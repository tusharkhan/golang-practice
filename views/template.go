package views

import (
	"bytes"
	"course/context"
	"course/models"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/csrf"
)

type PublicError interface {
	Public() string
}

type Template struct {
	htmlTemplate *template.Template
}

func ParseFS(fs fs.FS, pattern ...string) (Template, error) {
	htmlTpl := template.New(filepath.Base(pattern[0]))

	htmlTpl = htmlTpl.Funcs(
		template.FuncMap{
			"csrf": func() template.HTML {
				return `<! -- tag and ends with a -->`
			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("not implemented")
			},
			"errors": func() []string {
				return nil
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

func (t Template) Execute(writer http.ResponseWriter, request *http.Request, data interface{}, err ...error) {
	htmlTemplate, templateCloneError := t.htmlTemplate.Clone()

	if templateCloneError != nil {
		log.Printf("template cloning error %v", templateCloneError)
		http.Error(writer, "Error in rendering template", http.StatusInternalServerError)
		return
	}
	var errorMessages []string = PrintErrorMessages(err...)
	htmlTemplate = htmlTemplate.Funcs(
		template.FuncMap{
			"csrf": func() template.HTML {
				return csrf.TemplateField(request)
			},
			"currentUser": func() *models.User {
				return context.User(request.Context())
			},
			"errors": func() []string {
				return errorMessages
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

func PrintErrorMessages(err ...error) []string {
	var errorMessages []string
	fmt.Println(err)
	for _, message := range err {
		if message != nil {
			var pubError PublicError
			if errors.As(message, &pubError) {
				errorMessages = append(errorMessages, pubError.Public())
			} else {
				errorMessages = append(errorMessages, "Something went wrong")
			}
		}
	}

	return errorMessages
}
