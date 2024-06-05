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
	Render(app, w, r, "home.page.gohtml", nil)
}

func LoginPage(app *app.App, w http.ResponseWriter, r *http.Request) {
	Render(app, w, r, "login.page.gohtml", nil)
}

func LoginHandler(app *app.App, w http.ResponseWriter, r *http.Request) {
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

	valid, err := app.Models.User.PasswordMatches(email, password)
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

	http.Redirect(w, r, "/user/plans", http.StatusSeeOther)
}

func LogoutHandler(app *app.App, w http.ResponseWriter, r *http.Request) {
	userId := app.Session.PopString(r.Context(), UserInfoId)
	app.Session.Destroy(r.Context())
	app.Session.RenewToken(r.Context())

	app.InfoLog.Printf("%s logged out successfully", userId)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RegisterPage(app *app.App, w http.ResponseWriter, r *http.Request) {
	Render(app, w, r, "register.page.gohtml", nil)
}

func RegisterHandler(app *app.App, w http.ResponseWriter, r *http.Request) {
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

	app.Session.Put(r.Context(), TDFlash, "Check your email to activate your account")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ActivateAccountHandler(app *app.App, w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	appURL := fmt.Sprintf("%s%s", os.Getenv("BASE_URL"), uri)
	okay := utils.VerifyToken(appURL)

	if !okay {
		err := errors.New("invalid email token")
		app.ErrorLog.Printf("unable to activate: %s", err.Error())

		app.Session.Put(r.Context(), TDError, err.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.URL.Query().Get("email")
	user, err := app.Models.User.GetByEmail(email)
	if err != nil {
		app.ErrorLog.Printf("unable to activate: %s", err.Error())

		app.Session.Put(r.Context(), TDError, "no user found")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user.Active = true
	err = user.Update()
	if err != nil {
		app.ErrorLog.Printf("unable to activate: %s", err.Error())

		app.Session.Put(r.Context(), TDError, "unable to update user")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), TDFlash, "Account activated. You can now log in.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ForgotPasswordPage(app *app.App, w http.ResponseWriter, r *http.Request) {
	Render(app, w, r, "forgot-password.page.gohtml", nil)
}

func ForgotPasswordHandler(app *app.App, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Printf("unable to reset password: %s", err.Error())

		app.Session.Put(r.Context(), TDError, "please try again")
		http.Redirect(w, r, "/forgot-password", http.StatusSeeOther)
		return
	}

	email := r.Form.Get("email")

	exists := app.Models.User.UserExists(email)
	if !exists {
		app.ErrorLog.Println("unable to reset password: user does not exists")

		app.Session.Put(r.Context(), TDError, "invalid email")
		http.Redirect(w, r, "/forgot-password", http.StatusSeeOther)
		return
	}

	url := fmt.Sprintf("%s/reset-password?email=%s", os.Getenv("BASE_URL"), email)
	signedURL := utils.GenerateTokenFromString(url)
	app.InfoLog.Printf("signed url for %s is %s", email, signedURL)

	msg := emailer.Message{
		To:       email,
		Subject:  "Reset your password",
		Template: "reset-password-email",
		Data:     template.HTML(signedURL),
	}
	app.SendEmail(msg)

	app.Session.Put(r.Context(), TDFlash, "Check your email to reset password")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ResetPasswordPage(app *app.App, w http.ResponseWriter, r *http.Request) {
	Render(app, w, r, "reset-password.page.gohtml", nil)
}

func ResetPasswordHandler(app *app.App, w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	appURL := fmt.Sprintf("%s%s", os.Getenv("BASE_URL"), uri)
	app.InfoLog.Printf("app url received: %s", appURL)

	okay := utils.VerifyToken(appURL)

	if !okay {
		err := errors.New("invalid email token")
		app.ErrorLog.Printf("unable to reset password: %s", err.Error())

		app.Session.Put(r.Context(), TDError, err.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Printf("unable to reset password: %s", err.Error())

		app.Session.Put(r.Context(), TDError, "please try again")
		http.Redirect(w, r, "/forgot-password", http.StatusSeeOther)
		return
	}

	password := r.Form.Get("password")
	verifiedPassword := r.Form.Get("verify-password")

	if !utils.ValidatePassword(password) {
		err := errors.New("password is invalid")
		app.ErrorLog.Printf("unable to reset password: %s", err.Error())

		app.Session.Put(r.Context(), TDError, err.Error())
		http.Redirect(w, r, "/forgot-password", http.StatusSeeOther)
		return
	}

	if !strings.EqualFold(password, verifiedPassword) {
		err := errors.New("password does not match")
		app.ErrorLog.Printf("unable to reset password: %s", err.Error())

		app.Session.Put(r.Context(), TDError, err.Error())
		http.Redirect(w, r, "/forgot-password", http.StatusSeeOther)
		return
	}

	email := r.URL.Query().Get("email")
	user, err := app.Models.User.GetByEmail(email)
	if err != nil {
		app.ErrorLog.Printf("unable to reset password: %s", err.Error())

		app.Session.Put(r.Context(), TDError, "no user found")
		http.Redirect(w, r, "/forgot-password", http.StatusSeeOther)
		return
	}

	err = user.ResetPassword(password)
	if err != nil {
		app.ErrorLog.Printf("unable to reset password: %s", err.Error())

		app.Session.Put(r.Context(), TDError, "no user found")
		http.Redirect(w, r, "/forgot-password", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), TDFlash, "Password reset. You can now log in.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func PlansPage(app *app.App, w http.ResponseWriter, r *http.Request) {
	plans, err := app.Models.Plan.GetAll()
	if err != nil {
		app.ErrorLog.Printf("unable to chose subscription: %s", err.Error())
		app.Session.Put(r.Context(), TDError, "something went wrong")
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	dataMap := make(map[string]any)
	dataMap["plans"] = plans

	Render(app, w, r, "plans.page.gohtml", &TemplateData{
		Data: dataMap,
	})
}

func SubscriptionHandler(app *app.App, w http.ResponseWriter, r *http.Request) {
	planID := r.URL.Query().Get("id")

	plan, err := app.Models.Plan.GetOne(planID)
	if err != nil {
		app.ErrorLog.Printf("unable to subscribe: %s", err.Error())
		app.Session.Put(r.Context(), TDError, "no plans found")
		http.Redirect(w, r, "/user/plans", http.StatusSeeOther)

		return
	}

	user, ok := app.Session.Get(r.Context(), UserInfo).(database.User)
	if !ok {
		err := errors.New("no user found. please log in")
		app.ErrorLog.Printf("unable to subscribe: %s", err.Error())
		app.Session.Put(r.Context(), TDError, err.Error())
		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return
	}

	app.WG.Add(1)
	go utils.SendInvoice(app, user, plan)

	app.WG.Add(1)
	go utils.SendManual(app, user, plan)

	err = app.Models.Plan.SubscribeUserToPlan(user.ID, plan.ID)
	if err != nil {
		app.ErrorLog.Printf("unable to subscribe: %s", err.Error())
		app.Session.Put(r.Context(), TDError, "something went wrong. please try again")
		http.Redirect(w, r, "/user/plans", http.StatusSeeOther)

		return
	}

	updatedUser, err := app.Models.User.GetOne(user.ID)
	if err != nil {
		err := errors.New("error getting user from database")
		app.ErrorLog.Printf("unable to subscribe: %s", err.Error())
		app.Session.Put(r.Context(), TDError, err.Error())
		http.Redirect(w, r, "/user/plans", http.StatusSeeOther)

		return
	}

	app.Session.Put(r.Context(), UserInfo, updatedUser)

	app.Session.Put(r.Context(), TDFlash, "Subscribed!")
	http.Redirect(w, r, "/user/plans", http.StatusSeeOther)
}
