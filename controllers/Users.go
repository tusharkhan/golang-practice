package controller

import (
	"course/helper"
	"course/models"
	"encoding/json"
	"net/http"
)

type Users struct {
	Template struct {
		New Template
	}
}

func (u Users) New(writer http.ResponseWriter, request *http.Request) {
	u.Template.New.Execute(writer, nil)
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
