package main

import (
	"fmt"
	"github.com/pandaemoniumplaza/goroutines/subscription_service/data"
	"html/template"
	"net/http"
)

const (
	homePageTemplate     = "/home.page.gohtml"
	loginPageTemplate    = "/login.page.gohtml"
	registerPageTemplate = "/register.page.gohtml"
)

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, homePageTemplate, nil)
}

func (a *App) LoginPage(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, loginPageTemplate, nil)
}

func (a *App) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	_ = a.Session.RenewToken(r.Context())

	if err := r.ParseForm(); err != nil {
		a.ErrorLog.Println("failed to parse form")
		a.Session.Put(r.Context(), "error", "Something went wrong")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var (
		email = r.Form.Get("email")
		pwd   = r.Form.Get("password")
	)

	user, err := a.Models.User.GetByEmail(email)
	if err != nil {
		a.ErrorLog.Println("failed to query user", err)
		a.Session.Put(r.Context(), "error", "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// TODO: Check whether user is activated

	b, err := user.PasswordMatches(pwd)
	if err != nil {
		a.ErrorLog.Println("failed to verify password", err)
		a.Session.Put(r.Context(), "error", "Something went wrong")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !b {
		a.ErrorLog.Println("password doesn't match")

		msg := Message{
			To:      email,
			Subject: "Login attempt",
			Data:    "Failed login attempt - invalid password",
		}

		a.sendEmail(msg)

		a.Session.Put(r.Context(), "error", "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	a.Session.Put(r.Context(), "userID", user.ID)
	a.Session.Put(r.Context(), "flash", "Successful login!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *App) RegisterPage(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, registerPageTemplate, nil)
}

func (a *App) PostRegisterPage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		a.ErrorLog.Println("failed to parse form")
		a.Session.Put(r.Context(), "error", "Something went wrong")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// TODO: Check whether user is already registered

	u := data.User{
		Email:     r.Form.Get("email"),
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Password:  r.Form.Get("password"), // TODO: Generate from hash?
		Active:    0,
		IsAdmin:   0,
	}

	if _, err := u.Insert(u); err != nil {
		a.ErrorLog.Println("failed to create user")
		a.Session.Put(r.Context(), "error", "Something went wrong")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	a.InfoLog.Println("user has been save to db")

	signedUrl := GenerateTokenFromString(fmt.Sprintf("http://localhost:8080/activate?email=%s", u.Email))

	a.InfoLog.Println("signed url", signedUrl)

	a.InfoLog.Println("sending activation email to user")

	msg := Message{
		To:       u.Email,
		Subject:  "Active your account",
		Template: "confirmation-email",
		Data:     template.HTML(signedUrl),
	}

	a.sendEmail(msg)

	a.Session.Put(r.Context(), "flash", "Confirmation email sent! Check your email.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (a *App) ActiveAccount(w http.ResponseWriter, r *http.Request) {
	a.InfoLog.Println("received activation request")

	url := r.RequestURI
	testURL := fmt.Sprintf("http://localhost:8080%s", url)
	okay := VerifyToken(testURL)

	if !okay {
		a.ErrorLog.Println("failed verify token from activation link")
		a.Session.Put(r.Context(), "error", "Invalid token")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	a.InfoLog.Println("url verification successful")

	u, err := a.Models.User.GetByEmail(r.URL.Query().Get("email"))
	if err != nil {
		a.ErrorLog.Println("failed to find user by email")
		a.Session.Put(r.Context(), "error", "No user found")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	u.Active = 1
	if err = u.Update(); err != nil {
		a.ErrorLog.Println("failed to change user status")

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	a.InfoLog.Println("user has been activated")

	a.Session.Put(r.Context(), "flash", "Account has been activated successfully.")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
