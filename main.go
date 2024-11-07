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
	// database
	database, databaseError := helper.ConnectDatabase()

	if databaseError != nil {
		panic(databaseError)
	}

	defer database.Close()

	// controllers and services
	var userC controller.Users = controller.Users{
		UserService: &models.UserService{
			DB: database,
		},
		SessionService: &models.SessionService{
			DB: database,
		},
		PasswordResetService: &models.PasswordResetService{
			DB: database,
		},
	}

	var galleryController controller.Galleries = controller.Galleries{
		GalleryService: &models.GalleryService{
			DB: database,
		},
	}

	var umr controller.UserMiddleware = controller.UserMiddleware{
		SessionService: userC.SessionService,
	}

	// middleware
	var csrfString string = "007c4bf36082fc848409e97538568a9f2"

	csrfFunc := csrf.Protect([]byte(csrfString), csrf.Secure(false))

	router := chi.NewRouter()

	router.Use(csrfFunc)
	router.Use(umr.SetUser)

	homeTemplate := views.Must(views.ParseFS(templates.FS, "home.gohtml", "layout.gohtml"))
	contactTemplate := views.Must(views.ParseFS(templates.FS, "contact.gohtml", "layout.gohtml"))
	faqTemplate := views.Must(views.ParseFS(templates.FS, "faq.gohtml", "layout.gohtml"))
	createFAQ := views.Must(views.ParseFS(templates.FS, "faqCreate.gohtml", "layout.gohtml"))
	userC.Template.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml"))
	loginGet := views.Must(views.ParseFS(templates.FS, "signin.gohtml"))
	userC.Template.ForgetPasswordRequestForm = views.Must(views.ParseFS(templates.FS, "requestForgetPassword.gohtml"))
	userC.Template.ForgetPasswordSuccess = views.Must(views.ParseFS(templates.FS, "requestForgetPasswordSuccess.gohtml"))
	userC.Template.ChangePasswordView = views.Must(views.ParseFS(templates.FS, "requestForgetPasswordChange.gohtml"))

	galleryController.Template.New = views.Must(views.ParseFS(templates.FS, "galleryCreate.gohtml", "layout.gohtml"))

	router.Use(middleware.Logger)

	router.Get("/signup", userC.New)
	router.Post("/signup", userC.Create)

	router.Get("/signin", controller.StaticHandler(loginGet))
	router.Post("/signin", userC.LoginPOST)
	router.Post("/signout", userC.SignOut)

	router.Get("/forget-password", userC.ForgetPasswordRequestForm)
	router.Post("/forget-password", userC.ForgetPasswordRequest)

	router.Get("/forget-password-success", userC.ForgetPasswordRequestSuccess)
	router.Get("/forget", userC.ChangePasswordView)
	router.Post("/change-pass", userC.ChangePassword)

	router.Route("/user/me", func(r chi.Router) {
		r.Use(umr.RequireUser)
		r.Get("/", userC.CurrentUser)
	})

	router.Get("/", controller.StaticHandler(homeTemplate))

	router.Get("/contact", controller.StaticHandler(contactTemplate))

	router.Route("/faq", func(r chi.Router) {
		r.Get("/", controller.FaqHandler(faqTemplate))
		r.Get("/create", controller.FaqHandler(createFAQ))
		r.Post("/create", controller.CreateFAQ)
	})

	router.Route("/gallery", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(umr.RequireUser)
			r.Get("/", galleryController.New)
			r.Post("/create", galleryController.Create)
			r.Get("/delete/{id}", galleryController.Delete)
			r.Get("/{id}/edit", galleryController.Edit)
			r.Post("/{id}/edit", galleryController.EditPost)
		})
	})

	router.NotFound(notFoundHandler)
	http.ListenAndServe(":8080", router)
}
