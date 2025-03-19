package main

import "net/http"

func (a *App) SessionLoad(next http.Handler) http.Handler {
	return a.Session.LoadAndSave(next)
}

func (a *App) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.Session.Exists(r.Context(), "userID") {
			a.Session.Put(r.Context(), "error", "Log in first!")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		}

		next.ServeHTTP(w, r)
	})
}
