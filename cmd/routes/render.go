package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/app"
)

const TDError = "error"
const TDWarning = "warning"
const TDFlash = "flash"

const UserInfoId = "user_id"
const UserInfo = "user"

type TemplateData struct {
	Data          map[string]any
	Flash         string
	Warning       string
	Error         string
	Authenticated bool
	Now           time.Time
}

func render(app *app.App, w http.ResponseWriter, r *http.Request, tName string, tData *TemplateData) {
	partials := []string{
		fmt.Sprintf("%s/base.layout.gohtml", PathToTemplates),
		fmt.Sprintf("%s/header.partial.gohtml", PathToTemplates),
		fmt.Sprintf("%s/alerts.partial.gohtml", PathToTemplates),
		fmt.Sprintf("%s/navbar.partial.gohtml", PathToTemplates),
		fmt.Sprintf("%s/footer.partial.gohtml", PathToTemplates),
	}

	templates := []string{fmt.Sprintf("%s/%s", PathToTemplates, tName)}
	templates = append(templates, partials...)

	if tData == nil {
		tData = &TemplateData{}
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		app.ErrorLog.Printf("Error in parsing template files: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tData = addDefaultTemplateData(app, r, tData)
	err = tmpl.Execute(w, tData)
	if err != nil {
		app.ErrorLog.Printf("Error in executing template files: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addDefaultTemplateData(app *app.App, r *http.Request, td *TemplateData) *TemplateData {
	td.Flash = app.Session.PopString(r.Context(), TDFlash)
	td.Warning = app.Session.PopString(r.Context(), TDWarning)
	td.Error = app.Session.PopString(r.Context(), TDError)

	if td.Data == nil {
		td.Data = make(map[string]any)
	}
	td.Now = time.Now()
	td.Authenticated = isAuthenticated(app, r)

	return td
}

func isAuthenticated(app *app.App, r *http.Request) bool {
	return app.Session.Exists(r.Context(), UserInfoId)
}
