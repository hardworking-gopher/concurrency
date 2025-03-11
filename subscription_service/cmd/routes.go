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
	mux.Get("/register", a.RegisterPage)
	mux.Post("/register", a.PostRegisterPage)
	mux.Post("/activate-account", a.ActiveAccount)

	mux.Get("/test-email", func(writer http.ResponseWriter, request *http.Request) {
		m := Mailer{
			Domain:      "localhost",
			Host:        "localhost",
			Port:        1025,
			Username:    "",
			Password:    "",
			Encryption:  "none",
			FromAddress: "info@mycompany.com",
			FromName:    "info",
			Wait:        nil,
			ErrorChan:   make(chan error),
			DoneChan:    nil,
		}

		msg := Message{
			To:      "me@here.com",
			Subject: "Test email",
			Data:    "Hello, world!",
		}

		m.sendMail(msg, make(chan error))
	})

	return mux
}
