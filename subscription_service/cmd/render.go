package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

const (
	pathToTemplate = "subscription_service/templates"

	alertsPartialTemplate = "alerts.partial.gohtml"
	footerPartialTemplate = "footer.partial.gohtml"
	headerPartialTemplate = "header.partial.gohtml"
	navbarPartialTemplate = "navbar.partial.gohtml"
	baseLayoutTemplate    = "base.layout.gohtml"

	errorFailedToRender = "failed to render page"
)

type TemplateData struct {
	StringMap     map[string]string
	IntMap        map[string]int
	FloatMap      map[string]float64
	Data          map[string]any
	Flash         string
	Warning       string
	Error         string
	Authenticated bool
	Now           time.Time
	//User *data.User
}

func (a *App) render(w http.ResponseWriter, r *http.Request, t string, data *TemplateData) {
	templates := []string{
		fmt.Sprintf("%s/%s", pathToTemplate, t),

		fmt.Sprintf("%s/%s", pathToTemplate, headerPartialTemplate),
		fmt.Sprintf("%s/%s", pathToTemplate, baseLayoutTemplate),
		fmt.Sprintf("%s/%s", pathToTemplate, navbarPartialTemplate),
		fmt.Sprintf("%s/%s", pathToTemplate, alertsPartialTemplate),
		fmt.Sprintf("%s/%s", pathToTemplate, footerPartialTemplate),
	}

	if data == nil {
		data = &TemplateData{}
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		a.ErrorLog.Println(err)
		http.Error(w, errorFailedToRender, http.StatusInternalServerError)

		return
	}

	if err = tmpl.Execute(w, a.AppDefaultData(data, r)); err != nil {
		a.ErrorLog.Println(err)
		http.Error(w, errorFailedToRender, http.StatusInternalServerError)

		return
	}
}

func (a *App) AppDefaultData(data *TemplateData, r *http.Request) *TemplateData {
	data.Flash = a.Session.PopString(r.Context(), "flash")
	data.Warning = a.Session.PopString(r.Context(), "warning")
	data.Error = a.Session.PopString(r.Context(), "error")
	if a.IsAuthenticated(r) {
		data.Authenticated = true
	}
	data.Now = time.Now()

	return data
}

func (a *App) IsAuthenticated(r *http.Request) bool {
	return a.Session.Exists(r.Context(), "userID")
}
