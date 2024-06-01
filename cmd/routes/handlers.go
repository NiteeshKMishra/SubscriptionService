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
	app.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Printf("unable to login: %s", err.Error())

		app.Session.Put(r.Context(), TDError, "invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email, password := r.Form.Get("email"), r.Form.Get("password")

	user, err := app.Models.User.GetByEmail(email)
	if err != nil {
		app.ErrorLog.Printf("unable to login: %s", err.Error())

		app.Session.Put(r.Context(), TDError, "invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	valid, err := user.PasswordMatches(password)
	if err != nil {
		app.ErrorLog.Printf("unable to login: %s", err.Error())

		app.Session.Put(r.Context(), TDError, "invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !valid {
		app.ErrorLog.Println("password is not valid")

		app.Session.Put(r.Context(), TDError, "invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), UserInfoId, user.ID)
	app.Session.Put(r.Context(), UserInfo, user)
	app.Session.Put(r.Context(), TDFlash, "Logged In")

	app.InfoLog.Printf("%s logged in successfully", user.ID)

	http.Redirect(w, r, "/", http.StatusFound)
}

func LogoutPage(app *app.App, w http.ResponseWriter, r *http.Request) {
	userId := app.Session.PopString(r.Context(), UserInfoId)
	app.Session.Destroy(r.Context())
	app.Session.RenewToken(r.Context())

	app.InfoLog.Printf("%s logged out successfully", userId)

	http.Redirect(w, r, "/", http.StatusFound)
}

func RegisterPage(app *app.App, w http.ResponseWriter, r *http.Request) {
	render(app, w, r, "register.page.gohtml", nil)
}

func PostRegisterPage(app *app.App, w http.ResponseWriter, r *http.Request) {

}

func ActivateAccount(app *app.App, w http.ResponseWriter, r *http.Request) {

}
