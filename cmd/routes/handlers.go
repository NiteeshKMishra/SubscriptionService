package routes

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/app"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/database"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/emailer"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/utils"
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
		app.ErrorLog.Println("unable to login: password is not valid")
		msg := emailer.Message{
			To:      email,
			Subject: "Failed to log in",
			Data:    "password is not valid",
		}
		app.SendEmail(msg)

		app.Session.Put(r.Context(), TDError, "invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !user.Active {
		err := errors.New("user is not active")
		app.ErrorLog.Printf("unable to login: %s", err.Error())
		url := fmt.Sprintf("%s/activate?email=%s", os.Getenv("BASE_URL"), user.Email)
		signedURL := utils.GenerateTokenFromString(url)
		app.InfoLog.Printf("signed url for %s is %s", user.ID, signedURL)

		msg := emailer.Message{
			To:       user.Email,
			Subject:  "Activate your account",
			Template: "confirmation-email",
			Data:     template.HTML(signedURL),
		}
		app.SendEmail(msg)

		app.Session.Put(r.Context(), TDError, "Check your email to activate account, before login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), UserInfoId, user.ID)
	app.Session.Put(r.Context(), UserInfo, user)
	app.Session.Put(r.Context(), TDFlash, "Logged In")

	app.InfoLog.Printf("%s logged in successfully", user.ID)
	msg := emailer.Message{
		To:      email,
		Subject: "Login success",
		Data:    "Successfully log in",
	}
	app.SendEmail(msg)

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
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Printf("unable to register: %s", err.Error())

		app.Session.Put(r.Context(), TDError, "please try again")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	verifiedPassword := r.Form.Get("verify-password")
	firstName := r.Form.Get("first-name")
	lastName := r.Form.Get("last-name")

	if !utils.ValidateEmail(email) {
		err := errors.New("email is invalid")
		app.ErrorLog.Printf("unable to register: %s", err.Error())

		app.Session.Put(r.Context(), TDError, err.Error())
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	if !utils.ValidatePassword(password) {
		err := errors.New("password is invalid")
		app.ErrorLog.Printf("unable to register: %s", err.Error())

		app.Session.Put(r.Context(), TDError, err.Error())
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	if !strings.EqualFold(password, verifiedPassword) {
		err := errors.New("password does not match")
		app.ErrorLog.Printf("unable to register: %s", err.Error())

		app.Session.Put(r.Context(), TDError, err.Error())
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	if !utils.ValidateName(firstName) || !utils.ValidateName(lastName) {
		err := errors.New("name is not valid")
		app.ErrorLog.Printf("unable to register: %s", err.Error())

		app.Session.Put(r.Context(), TDError, err.Error())
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	user := database.User{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		IsAdmin:   false,
		Active:    false,
	}

	exists := user.UserExists(email)
	if exists {
		err := errors.New("user already exists. please login")
		app.ErrorLog.Printf("unable to register: %s", err.Error())

		app.Session.Put(r.Context(), TDError, err.Error())
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user.ID, err = user.Insert()
	if err != nil {
		app.ErrorLog.Printf("unable to register: %s", err.Error())

		app.Session.Put(r.Context(), TDError, "please try again")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// send an activation email
	url := fmt.Sprintf("%s/activate?email=%s", os.Getenv("BASE_URL"), user.Email)
	signedURL := utils.GenerateTokenFromString(url)
	app.InfoLog.Printf("signed url for %s is %s", user.ID, signedURL)

	msg := emailer.Message{
		To:       user.Email,
		Subject:  "Activate your account",
		Template: "confirmation-email",
		Data:     template.HTML(signedURL),
	}

	app.SendEmail(msg)

	app.Session.Put(r.Context(), "flash", "Check your email to activate your account")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ActivateAccount(app *app.App, w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	appURL := fmt.Sprintf("%s%s", os.Getenv("BASE_URL"), uri)
	okay := utils.VerifyToken(appURL)

	if !okay {
		err := errors.New("invalid email token")
		app.ErrorLog.Printf("unable to activate: %s", err.Error())

		app.Session.Put(r.Context(), "error", err.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.URL.Query().Get("email")
	user, err := app.Models.User.GetByEmail(email)
	if err != nil {
		app.ErrorLog.Printf("unable to activate: %s", err.Error())

		app.Session.Put(r.Context(), "error", "no user found")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user.Active = true
	err = user.Update()
	if err != nil {
		app.ErrorLog.Printf("unable to activate: %s", err.Error())

		app.Session.Put(r.Context(), "error", "unable to update user")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), "flash", "Account activated. You can now log in.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
