package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/app"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/database"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/utils"
)

const TDError = "error"
const TDWarning = "warning"
const TDFlash = "flash"

const UserInfoId = "user_id"
const UserInfo = "user"

func render(app *app.App, w http.ResponseWriter, r *http.Request, tName string, tData *TemplateData) {
	partials := []string{
		fmt.Sprintf("%s/base.layout.gohtml", utils.PathToTemplates),
		fmt.Sprintf("%s/header.partial.gohtml", utils.PathToTemplates),
		fmt.Sprintf("%s/alerts.partial.gohtml", utils.PathToTemplates),
		fmt.Sprintf("%s/navbar.partial.gohtml", utils.PathToTemplates),
		fmt.Sprintf("%s/footer.partial.gohtml", utils.PathToTemplates),
	}

	templates := []string{fmt.Sprintf("%s/%s", utils.PathToTemplates, tName)}
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
	if isAuthenticated(app, r) {
		td.Authenticated = true
		user, ok := app.Session.Get(r.Context(), UserInfo).(database.User)
		if !ok {
			app.ErrorLog.Println("can't get user from session")
		} else {
			td.User = &user
		}
	}

	return td
}

func isAuthenticated(app *app.App, r *http.Request) bool {
	return app.Session.Exists(r.Context(), UserInfoId)
}
