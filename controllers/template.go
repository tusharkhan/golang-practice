package controller

import "net/http"

type Template interface {
	Execute(writer http.ResponseWriter, request *http.Request, data interface{})
}
