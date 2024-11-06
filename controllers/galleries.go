package controller

import (
	"course/context"
	"course/models"
	"fmt"
	"net/http"
)

type Galleries struct {
	Template struct {
		New Template
	}

	GalleryService *models.GalleryService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	g.Template.New.Execute(w, r, nil, nil)
}

func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var userId int = context.User(r.Context()).ID
	parseError := r.ParseForm()

	if parseError != nil {
		panic(parseError)
	}
	var title string = r.FormValue("title")

	createGallery, createGalleryError := g.GalleryService.Create(title, userId)

	if createGalleryError != nil {
		g.Template.New.Execute(w, r, nil, createGalleryError)
		return
	}

	fmt.Fprintf(w, "%v", createGallery)
}
