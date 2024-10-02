package models

import (
	"course/helper"
	"database/sql"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) CreateUser(name, email, password string) (*User, error) {
	var emailLower string = strings.ToLower(email)

	pass, hashError := helper.HashString(password)

	if hashError != nil {
		return nil, hashError
	}

	insertedRow := us.DB.QueryRow("INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id;", name, emailLower, pass)

	var use User = User{
		Name:     name,
		Email:    emailLower,
		Password: pass,
	}

	insertError := insertedRow.Scan(&use.ID)

	if insertError != nil {
		return nil, insertError
	}

	us.DB.Close()

	return &use, nil
}
