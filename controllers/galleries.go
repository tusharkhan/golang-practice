package controller

import (
	"course/context"
	"course/helper"
	"course/models"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Galleries struct {
	Template struct {
		New  Template
		Show Template
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

	createdGallery, createGalleryError := g.GalleryService.Create(title, userId)

	if createGalleryError != nil {
		g.Template.New.Execute(w, r, nil, createGalleryError)
	}

	multipartParsingError := r.ParseMultipartForm(5 << 20)

	if multipartParsingError != nil {
		g.Template.New.Execute(w, r, nil, models.MaxImageSizeError)
	}

	fileHeaders := r.MultipartForm.File["gaslleryImages"]
	var uploadedImageName []models.GalleryImages
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		fmt.Println(file)
		if err != nil {
			g.Template.New.Execute(w, r, nil, models.InternalServerError)
		}
		defer file.Close()

		uploadedImage, imageUploadError := g.GalleryService.UploadImage(createdGallery.ID, fileHeader, file)

		if imageUploadError != nil {
			g.Template.New.Execute(w, r, nil, imageUploadError)
		}
		var img models.GalleryImages = models.GalleryImages{
			GalleryId:    uploadedImage.GalleryId,
			FileSize:     uploadedImage.FileSize,
			RealName:     uploadedImage.RealName,
			GenerateName: uploadedImage.GenerateName,
			CreatedAt:    uploadedImage.CreatedAt,
		}
		uploadedImageName = append(uploadedImageName, img)
	}

	inserrtError := g.GalleryService.InsertImage(uploadedImageName)

	if inserrtError != nil {
		g.Template.New.Execute(w, r, nil, inserrtError)
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

func (g Galleries) Show(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Id     int
		Title  string
		Images []models.Image
	}

	id, urlParamError := strconv.Atoi(chi.URLParam(r, "id"))

	if urlParamError != nil {
		http.Error(w, "Invalide url param", http.StatusInternalServerError)
	}

	singleGallery, singleGalleryError := g.GalleryService.Show(id)

	if singleGalleryError != nil {
		g.Template.New.Execute(w, r, nil, singleGalleryError)
	}

	galleryImages, galleryImagesError := g.GalleryService.Images(id)

	if galleryImagesError != nil {
		g.Template.New.Execute(w, r, nil, singleGalleryError)
	}

	data.Id = singleGallery.ID
	data.Title = singleGallery.Title
	data.Images = galleryImages

	g.Template.Show.Execute(w, r, data, nil)
}

func (g Galleries) RenderImage(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	galleryid, urlParamError := strconv.Atoi(chi.URLParam(r, "galleryid"))

	if urlParamError != nil {
		http.Error(w, "Invalide url param", http.StatusInternalServerError)
	}

	filepa := g.GalleryService.GalleryDire(galleryid)
	realPath := filepath.Join(filepa, filename)

	http.ServeFile(w, r, realPath)
}

func (g Galleries) DeleteImage(w http.ResponseWriter, r *http.Request) {
	galleryId := chi.URLParam(r, "galleryid")
	id, urlParamError := strconv.Atoi(galleryId)

	if urlParamError != nil {
		http.Error(w, "Invalide url param", http.StatusInternalServerError)
	}

	var filename string = chi.URLParam(r, "filename")

	var imageDeleteError error = g.GalleryService.RemoveImage(id, filename)

	if imageDeleteError != nil {
		http.Error(w, imageDeleteError.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}
