package controller

import "net/http"

type Template interface {
	Execute(writer http.ResponseWriter, data interface{})
}
