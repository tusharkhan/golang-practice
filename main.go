package main

import (
	controller "course/controllers"
	"course/helper"
	"course/models"
	"course/templates"
	"course/views"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
)

type User struct {
	Name string
	Age  int
}

func contactHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "<h1>Contact </h1>")
}

func notFoundHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(404)
	fmt.Fprint(writer, "<h1>404 Not found </h1>")
}

func main() {
	router := chi.NewRouter()

	homeTemplate := views.Must(views.ParseFS(templates.FS, "home.gohtml", "layout.gohtml"))
	contactTemplate := views.Must(views.ParseFS(templates.FS, "contact.gohtml", "layout.gohtml"))
	faqTemplate := views.Must(views.ParseFS(templates.FS, "faq.gohtml", "layout.gohtml"))
	createFAQ := views.Must(views.ParseFS(templates.FS, "faqCreate.gohtml", "layout.gohtml"))
	signup := views.Must(views.ParseFS(templates.FS, "signup.gohtml"))
	loginGet := views.Must(views.ParseFS(templates.FS, "signin.gohtml"))

	database, databaseError := helper.ConnectDatabase()

	if databaseError != nil {
		panic(databaseError)
	}

	defer database.Close()

	var userC controller.Users = controller.Users{
		UserService: &models.UserService{
			DB: database,
		},
		SessionService: &models.SessionService{
			DB: database,
		},
	}

	userC.Template.New = signup

	router.Use(middleware.Logger)

	router.Get("/signup", userC.New)
	router.Post("/signup", userC.Create)

	router.Get("/signin", controller.StaticHandler(loginGet))
	router.Post("/signin", userC.LoginPOST)

	router.Get("/", controller.StaticHandler(homeTemplate))

	router.Get("/contact", controller.StaticHandler(contactTemplate))

	router.Route("/faq", func(r chi.Router) {
		r.Get("/", controller.FaqHandler(faqTemplate))
		r.Get("/create", controller.FaqHandler(createFAQ))
		r.Post("/create", controller.CreateFAQ)
	})

	router.NotFound(notFoundHandler)

	var csrfString string = "007c4bf36082fc848409e97538568a9f2"

	csrfFunc := csrf.Protect([]byte(csrfString), csrf.Secure(false))

	http.ListenAndServe(":8080", csrfFunc(router))
}
