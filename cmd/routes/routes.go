package routes

import (
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/app"
)

var PathToTemplates = path.Join("./", "cmd", "templates")

func InitRoutes(app *app.App) http.Handler {
	mux := chi.NewRouter()
	//Middlewares
	mux.Use(middleware.Recoverer)
	mux.Use(func(h http.Handler) http.Handler {
		return app.Session.LoadAndSave(h)
	})

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
		PostLoginPage(app, w, r)
	})
	mux.Get("/register", func(w http.ResponseWriter, r *http.Request) {
		RegisterPage(app, w, r)
	})
	mux.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		PostRegisterPage(app, w, r)
	})
	mux.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		LogoutPage(app, w, r)
	})
	mux.Get("/activate", func(w http.ResponseWriter, r *http.Request) {
		ActivateAccount(app, w, r)
	})

	return mux
}
