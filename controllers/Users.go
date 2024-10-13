package controller

import (
	"course/helper"
	"course/models"
	"fmt"
	"net/http"
	"strings"
)

type Users struct {
	Template struct {
		New Template
	}

	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u Users) New(writer http.ResponseWriter, request *http.Request) {
	u.Template.New.Execute(writer, request, nil)
}

func (u Users) Create(writer http.ResponseWriter, request *http.Request) {
	parseError := request.ParseForm()

	if parseError != nil {
		panic(parseError)
	}

	var name string = request.FormValue("name")
	var email string = strings.ToLower(request.FormValue("email"))
	var password string = request.FormValue("password")

	createdUser, creatingError := u.UserService.CreateUser(name, email, password)

	if creatingError != nil {
		http.Error(writer, creatingError.Error(), http.StatusInternalServerError)
		return
	}

	userSession, sessionCreateError := u.SessionService.Create(createdUser.ID)

	if sessionCreateError != nil {
		http.Redirect(writer, request, "/signup", http.StatusFound)
		return
	}

	helper.SetNewCookie(writer, helper.CookieSession, userSession.Token)
	http.Redirect(writer, request, "/", http.StatusFound)
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
		fmt.Errorf("Error in Login %w", errorInLogin)
		http.Error(writer, "Error in Login User", http.StatusInternalServerError)
		return
	}

	passCheck := helper.CheckPassword(loggedInUser.Password, password)

	if !passCheck {
		http.Error(writer, "Invalid Credentials", http.StatusNotFound)
		return
	}

	createSession, sessionCreateError := u.SessionService.Create(loggedInUser.ID)

	if sessionCreateError != nil {
		fmt.Errorf("Error in Creating Token %w", errorInLogin)
		http.Error(writer, "Error in Creating Token", http.StatusInternalServerError)
		return
	}

	helper.SetNewCookie(writer, helper.CookieSession, createSession.Token)

	http.Redirect(writer, request, "/", http.StatusOK)
}

func (u Users) CurrentUser(writer http.ResponseWriter, request *http.Request) {
	cookie, tokenError := helper.ReadCookie(request, helper.CookieSession)

	if tokenError != nil {
		http.Redirect(writer, request, "/signin", http.StatusFound)
		return
	}

	user, userError := u.SessionService.User(cookie)

	if userError != nil {
		http.Redirect(writer, request, "/signin", http.StatusFound)
		return
	}

	fmt.Fprintf(writer, "Current User %s \n", user.Email)
}

func (u Users) Update(writer http.ResponseWriter, request *http.Request) {
	// var userId string = request.URL.Query().Get("id")

	// parseError := request.ParseForm()

	// if parseError != nil {
	// 	panic(parseError)
	// }

	// var email string = strings.ToLower(request.FormValue("email"))
	// var name string = request.FormValue("name")
	// var password string = request.FormValue("password")

	// update, updateError := u.UserService.UpdateUser(userId, name, email, password)

	// // TODO : update user function
	// if updateError != nil {
	// 	http.Redirect(writer, request, "/", http.StatusInternalServerError)
	// }

}
