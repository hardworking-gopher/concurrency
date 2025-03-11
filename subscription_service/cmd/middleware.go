package main

import "net/http"

func (a *App) SessionLoad(next http.Handler) http.Handler {
	return a.Session.LoadAndSave(next)
}
