package main

import (
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

	b, err := a.Models.User.PasswordMatches(pwd)
	if err != nil {
		a.ErrorLog.Println("failed to verify password", err)
		a.Session.Put(r.Context(), "error", "Something went wrong")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !b {
		a.ErrorLog.Println("password doesn't match")
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
	//a.render(w, r, loginPageTemplate, nil)
}

func (a *App) ActiveAccount(w http.ResponseWriter, r *http.Request) {
	//a.render(w, r, loginPageTemplate, nil)
}
