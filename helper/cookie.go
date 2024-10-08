package helper

import (
	"fmt"
	"net/http"
)

const (
	CookieSession = "session"
)

func SetNewCookie(w http.ResponseWriter, name, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	})
}

func ReadCookie(req *http.Request, name string) (string, error) {
	coo, err := req.Cookie(name)

	if err != nil {
		return "", fmt.Errorf("%s %w", name, err)
	}

	return coo.Value, nil
}
