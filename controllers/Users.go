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
		return
	}

	userSession, sessionCreateError := u.SessionService.Create(createdUser.ID)

	if sessionCreateError != nil {
		fmt.Println(sessionCreateError)
		http.Redirect(writer, request, "/signin", http.StatusFound)
	}

	cookie := http.Cookie{
		Name:     "Session",
		Value:    userSession.Token,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(writer, &cookie)
	http.Redirect(writer, request, "/", http.StatusFound)
	// writer.Header().Set("Content-Type", "application/json")

	// json.NewEncoder(writer).Encode(createdUser)
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

	cookie := http.Cookie{
		Name:     "Session",
		Value:    createSession.Token,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(writer, &cookie)
	http.Redirect(writer, request, "/", http.StatusOK)
}

func (u Users) CurrentUser(writer http.ResponseWriter, request *http.Request) {
	token, tokenError := request.Cookie("Session")

	if tokenError != nil {
		http.Redirect(writer, request, "/signin", http.StatusFound)
		return
	}

	user, userError := u.SessionService.User(token.Value)

	if userError != nil {
		http.Redirect(writer, request, "/signin", http.StatusFound)
		return
	}

	fmt.Fprintf(writer, "Current User %s \n", user.Email)
}
