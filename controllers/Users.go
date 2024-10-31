package controller

import (
	"course/context"
	"course/helper"
	"course/models"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Users struct {
	Template struct {
		New                       Template
		ForgetPasswordRequestForm Template
	}

	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
}

type UserMiddleware struct {
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

	loggedInUser, errorInLogin := u.UserService.Login(email, password)

	if errorInLogin != nil {
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
		http.Error(writer, "Error in Creating Token", http.StatusInternalServerError)
		return
	}

	helper.SetNewCookie(writer, helper.CookieSession, createSession.Token)

	http.Redirect(writer, request, "/", http.StatusFound)
}

func (u Users) CurrentUser(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	user := context.User(ctx)

	if user == nil {
		http.Redirect(writer, request, "/signin", http.StatusFound)
		return
	}

	fmt.Fprintf(writer, "Current User %s \n", user.Email)
}

func (u Users) SignOut(writer http.ResponseWriter, request *http.Request) {
	cookie, tokenError := helper.ReadCookie(request, helper.CookieSession)

	if tokenError != nil {
		http.Redirect(writer, request, "/signin", http.StatusFound)
		return
	}

	var tokenHash bool = u.SessionService.DestroySession(cookie)

	if !tokenHash {
		http.Error(writer, "Error in signout", http.StatusInternalServerError)
	} else {
		http.SetCookie(writer, &http.Cookie{
			Name:     helper.CookieSession,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
		http.Redirect(writer, request, "/", http.StatusFound)
	}
}

func (u Users) ForgetPasswordRequestForm(w http.ResponseWriter, r *http.Request) {
	u.Template.ForgetPasswordRequestForm.Execute(w, r, nil)
}

func (u Users) ForgetPasswordRequest(w http.ResponseWriter, r *http.Request) {
	parseError := r.ParseForm()

	if parseError != nil {
		panic(parseError)
	}

	var email string = r.FormValue("email")

	passReset, passwordResetError := u.PasswordResetService.Create(email)

	if passwordResetError != nil {
		http.Error(w, "Something went wrong creating token", http.StatusInternalServerError)
		return
	}

	userlWithToken := url.Values{
		"token": {passReset.Token},
	}
	var resetUrl string = helper.BaseURL(r) + "/forget" + userlWithToken.Encode()
	sendingMailError := u.EmailService.SendForgetPasswordEmail(email, resetUrl)

	if sendingMailError != nil {
		http.Error(w, "Something went wrong sending mail", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

func (umr UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, tokenReadError := helper.ReadCookie(r, helper.CookieSession)
		if tokenReadError != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, userError := umr.SessionService.User(token)
		if userError != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (umr UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())

		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
