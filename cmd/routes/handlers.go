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
	app.InfoLog.Println("home page")
	fmt.Fprint(w, "home page")
}
