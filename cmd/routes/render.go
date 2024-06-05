package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/app"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/constants"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/database"
)

const TDError = "error"
const TDWarning = "warning"
const TDFlash = "flash"

const UserInfoId = "user_id"
const UserInfo = "user"

func Render(app *app.App, w http.ResponseWriter, r *http.Request, tName string, tData *TemplateData) {
	partials := []string{
		fmt.Sprintf("%s/base.layout.gohtml", constants.PathToTemplates),
		fmt.Sprintf("%s/header.partial.gohtml", constants.PathToTemplates),
		fmt.Sprintf("%s/alerts.partial.gohtml", constants.PathToTemplates),
		fmt.Sprintf("%s/navbar.partial.gohtml", constants.PathToTemplates),
		fmt.Sprintf("%s/footer.partial.gohtml", constants.PathToTemplates),
	}

	templates := []string{fmt.Sprintf("%s/%s", constants.PathToTemplates, tName)}
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

	tData = AddDefaultTemplateData(app, r, tData)
	err = tmpl.Execute(w, tData)
	if err != nil {
		app.ErrorLog.Printf("Error in executing template files: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AddDefaultTemplateData(app *app.App, r *http.Request, td *TemplateData) *TemplateData {
	td.Flash = app.Session.PopString(r.Context(), TDFlash)
	td.Warning = app.Session.PopString(r.Context(), TDWarning)
	td.Error = app.Session.PopString(r.Context(), TDError)

	if td.Data == nil {
		td.Data = make(map[string]any)
	}
	if IsAuthenticated(app, r) {
		td.Authenticated = true
		user, ok := app.Session.Get(r.Context(), UserInfo).(database.User)
		if !ok {
			app.ErrorLog.Println("can't get user from session")
		} else {
			td.User = &user
		}
	}

	td.Now = time.Now()
	return td
}

func IsAuthenticated(app *app.App, r *http.Request) bool {
	return app.Session.Exists(r.Context(), UserInfoId)
}
