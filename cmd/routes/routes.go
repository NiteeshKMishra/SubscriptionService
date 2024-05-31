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

	return mux
}
