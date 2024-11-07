package controller

import (
	"course/context"
	"course/helper"
	"course/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Galleries struct {
	Template struct {
		New Template
	}

	GalleryService *models.GalleryService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	list, listError := g.GalleryService.List()

	if listError != nil {
		g.Template.New.Execute(w, r, nil, listError)
	}

	var data struct {
		ListGallery []models.Gallery
	}

	data.ListGallery = list

	g.Template.New.Execute(w, r, data, nil)
}

func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var userId int = context.User(r.Context()).ID
	parseError := r.ParseForm()

	if parseError != nil {
		panic(parseError)
	}
	var title string = r.FormValue("title")

	_, createGalleryError := g.GalleryService.Create(title, userId)

	if createGalleryError != nil {
		g.Template.New.Execute(w, r, nil, createGalleryError)
	}

	http.Redirect(w, r, "/gallery", http.StatusFound)
}

func (g Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	id, urlParamError := strconv.Atoi(chi.URLParam(r, "id"))

	if urlParamError != nil {
		http.Error(w, "Invalide url param", http.StatusInternalServerError)
	}

	deleteError := g.GalleryService.Delete(id)

	if deleteError != nil {
		http.Error(w, deleteError.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/gallery", http.StatusFound)
}

func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	id, urlParamError := strconv.Atoi(chi.URLParam(r, "id"))

	if urlParamError != nil {
		http.Error(w, "Invalide url param", http.StatusInternalServerError)
	}

	gallery, getGalleryError := g.GalleryService.Show(id)

	if getGalleryError != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	helper.ResponseJSON(w, r, gallery)
}

func (g Galleries) EditPost(w http.ResponseWriter, r *http.Request) {
	id, urlParamError := strconv.Atoi(chi.URLParam(r, "id"))

	if urlParamError != nil {
		http.Error(w, "Invalide url param", http.StatusInternalServerError)
	}

	parseError := r.ParseForm()

	if parseError != nil {
		panic(parseError)
	}

	var title string = r.FormValue("title")

	updateError := g.GalleryService.Update(&models.Gallery{
		ID:    id,
		Title: title,
	})

	if updateError != nil {
		http.Error(w, updateError.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/gallery", http.StatusFound)
}
