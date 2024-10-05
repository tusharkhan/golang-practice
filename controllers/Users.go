package controller

import (
	"course/helper"
	"course/models"
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/csrf"
)

type Users struct {
	Template struct {
		New Template
	}
}

func (u Users) New(writer http.ResponseWriter, request *http.Request) {
	var data struct {
		CsrfField template.HTML
	}

	data.CsrfField = csrf.TemplateField(request)

	u.Template.New.Execute(writer, data)
}

func (u Users) Create(writer http.ResponseWriter, request *http.Request) {
	parseError := request.ParseForm()

	if parseError != nil {
		panic(parseError)
	}

	var name string = request.FormValue("name")
	var email string = request.FormValue("email")
	var password string = request.FormValue("password")

	database, databaseError := helper.ConnectDatabase()

	if databaseError != nil {
		panic(databaseError)
	}

	defer database.Close()

	createUser := models.UserService{
		DB: database,
	}

	createdUser, creatingError := createUser.CreateUser(name, email, password)

	if creatingError != nil {
		http.Error(writer, creatingError.Error(), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")

	json.NewEncoder(writer).Encode(createdUser)
}

func (u Users) LoginPOST(writer http.ResponseWriter, request *http.Request) {
	parseError := request.ParseForm()

	if parseError != nil {
		panic(parseError)
	}

	var email string = strings.ToLower(request.FormValue("email"))
	var password string = request.FormValue("password")

	database, databaseConnectionError := helper.ConnectDatabase()

	if databaseConnectionError != nil {
		panic(databaseConnectionError)
	}

	defer database.Close()

	var loginUser models.UserService = models.UserService{
		DB: database,
	}

	loggedInUser, errorInLogin := loginUser.Login(email, password)

	if errorInLogin != nil {
		http.Redirect(writer, request, "/signin", http.StatusSeeOther)
	}

	passCheck := helper.CheckPassword(loggedInUser.Password, password)

	if !passCheck {
		http.Redirect(writer, request, "/signin", http.StatusSeeOther)
	}

	cookie := http.Cookie{
		Name:     "email",
		Value:    loggedInUser.Email,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(writer, &cookie)
}
