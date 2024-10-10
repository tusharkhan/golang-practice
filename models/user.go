package models

import (
	"course/helper"
	"database/sql"
	"strings"
	"time"
)

type User struct {
	ID         int
	Name       string
	Email      string
	Password   string
	Created_at string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) CreateUser(name, email, password string) (*User, error) {
	var emailLower string = strings.ToLower(email)

	pass, hashError := helper.HashString(password)

	created_at := time.Now().Format("2006-01-02 15:04:05")

	if hashError != nil {
		return nil, hashError
	}

	insertedRow := us.DB.QueryRow("INSERT INTO users (name, email, password, created_at) VALUES ($1, $2, $3, $4) RETURNING id;", name, emailLower, pass, created_at)

	var use User = User{
		Name:       name,
		Email:      emailLower,
		Password:   pass,
		Created_at: created_at,
	}

	insertError := insertedRow.Scan(&use.ID)

	if insertError != nil {
		return nil, insertError
	}

	return &use, nil
}

func (us *UserService) Login(email, password string) (*User, error) {
	var sql string = "SELECT id, name, email, password, created_at FROM users WHERE email=$1"
	email = strings.ToLower(email)

	row := us.DB.QueryRow(sql, email)
	var userFromQuery User = User{
		Email: email,
	}
	getDataError := row.Scan(&userFromQuery.ID, &userFromQuery.Name, &userFromQuery.Email, &userFromQuery.Password, &userFromQuery.Created_at)

	if getDataError != nil {
		return nil, getDataError
	}

	return &userFromQuery, nil
}

func (u *UserService) UpdateUser(id, name, email, password string) (*User, error) {
	var query string = "UPDATE users SET name=$1, email=$2, password=$3 WHERE id=$4"
	var updatedUser User = User{}
	pass, hashError := helper.HashString(password)

	if hashError != nil {
		return nil, hashError
	}

	row := u.DB.QueryRow(query, name, email, pass, id)

	erro := row.Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Email, &updatedUser.Password)

	if erro != nil {
		return nil, erro
	}

	return &updatedUser, nil
}
