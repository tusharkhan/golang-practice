package models

import "errors"

var (
	EmailAlreadyTaken = errors.New("models : Email already taken")
	ErrorNotFound     = errors.New("models : No record found")
)
