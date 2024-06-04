package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/app"
)

func InitRoutes(app *app.App) http.Handler {
	mux := chi.NewRouter()
	//Middlewares
	mux.Use(middleware.Recoverer)
	mux.Use(app.Session.LoadAndSave)

	//Routes
	mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		Health(app, w, r)
	})
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		HomePage(app, w, r)
	})
	mux.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		LoginPage(app, w, r)
	})
	mux.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		LoginHandler(app, w, r)
	})
	mux.Get("/register", func(w http.ResponseWriter, r *http.Request) {
		RegisterPage(app, w, r)
	})
	mux.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		RegisterHandler(app, w, r)
	})
	mux.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		LogoutHandler(app, w, r)
	})
	mux.Get("/forgot-password", func(w http.ResponseWriter, r *http.Request) {
		ForgotPasswordPage(app, w, r)
	})
	mux.Post("/forgot-password", func(w http.ResponseWriter, r *http.Request) {
		ForgotPasswordHandler(app, w, r)
	})
	mux.Get("/reset-password", func(w http.ResponseWriter, r *http.Request) {
		ResetPasswordPage(app, w, r)
	})
	mux.Post("/reset-password", func(w http.ResponseWriter, r *http.Request) {
		ResetPasswordHandler(app, w, r)
	})
	mux.Get("/activate", func(w http.ResponseWriter, r *http.Request) {
		ActivateAccountHandler(app, w, r)
	})

	mux.Mount("/user", authRouter(app))

	return mux
}

func authRouter(app *app.App) http.Handler {
	mux := chi.NewRouter()
	mux.Use(Auth(app))

	mux.Get("/plans", func(w http.ResponseWriter, r *http.Request) {
		PlansPage(app, w, r)
	})
	mux.Get("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		SubscriptionHandler(app, w, r)
	})

	return mux
}
