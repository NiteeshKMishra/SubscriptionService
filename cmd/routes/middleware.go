package routes

import (
	"net/http"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/app"
)

func Auth(app *app.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !app.Session.Exists(r.Context(), UserInfoId) {
				app.Session.Put(r.Context(), TDError, "Please log in to continue")
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
