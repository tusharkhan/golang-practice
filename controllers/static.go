package controller

import (
	"course/helper"
	"course/models"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

var faqModel models.FAQ = models.FAQ{}

func StaticHandler(temp Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		temp.Execute(w, r, nil)
	}
}

func FaqHandler(temp Template) http.HandlerFunc {
	fetAllData := faqModel.Get()

	return func(w http.ResponseWriter, r *http.Request) {
		temp.Execute(w, r, fetAllData)
	}
}

func ShwoFaq(temp Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var faqs []models.FAQ = faqModel.Get()
		id := chi.URLParam(r, "id")
		intId, parseError := strconv.ParseInt(id, 10, 4)
		fmt.Println(id, intId)
		if parseError != nil {
			panic(parseError)
		}

		temp.Execute(w, r, faqs[intId-1])
	}
}

func CreateFAQ(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()

	if err != nil {
		panic(err)
	}

	var faqInput models.FAQ = models.FAQ{
		Question: template.HTML(request.FormValue("question")),
		Answer:   template.HTML(request.FormValue("answer")),
		UserRating: models.Ratings{
			User:  request.FormValue("user"),
			Email: request.FormValue("email"),
			Image: request.FormValue("image"),
		},
	}

	openFile, openFileError := helper.OpenFile(models.FAQ.GetRealFilePath(faqModel))

	if openFileError != nil {
		panic(openFileError)
	}

	structToJson, jsonError := json.Marshal(faqInput)

	if jsonError != nil {
		panic(jsonError)
	}

	helper.WriteFile(&openFile, string(structToJson))
}

func SignupPage(temp Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		temp.Execute(w, r, nil)
	}
}
