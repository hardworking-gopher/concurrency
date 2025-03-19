package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (a *App) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(a.SessionLoad)

	mux.Get("/", a.Home)
	mux.Get("/login", a.LoginPage)
	mux.Post("/login", a.PostLoginPage)
	mux.Get("/logout", a.Logout)
	mux.Get("/register", a.RegisterPage)
	mux.Post("/register", a.PostRegisterPage)
	mux.Get("/activate", a.ActiveAccount)

	mux.Mount("/members", a.authRouter())

	return mux
}

func (a *App) authRouter() http.Handler {
	mux := chi.NewRouter()
	mux.Use(a.Auth)

	mux.Get("/plans", a.ChooseSubscription)
	mux.Get("/subscribe", a.Subscribe)

	return mux
}
