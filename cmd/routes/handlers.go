package routes

import (
	"fmt"
	"net/http"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/app"
)

func Health(app *app.App, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "app is running with status %d", http.StatusOK)
}

func HomePage(app *app.App, w http.ResponseWriter, r *http.Request) {
	render(app, w, r, "home.page.gohtml", nil)
}

func LoginPage(app *app.App, w http.ResponseWriter, r *http.Request) {
	render(app, w, r, "login.page.gohtml", nil)
}

func PostLoginPage(app *app.App, w http.ResponseWriter, r *http.Request) {

}

func LogoutPage(app *app.App, w http.ResponseWriter, r *http.Request) {

}

func RegisterPage(app *app.App, w http.ResponseWriter, r *http.Request) {
	render(app, w, r, "register.page.gohtml", nil)
}

func PostRegisterPage(app *app.App, w http.ResponseWriter, r *http.Request) {

}

func ActivateAccount(app *app.App, w http.ResponseWriter, r *http.Request) {

}
