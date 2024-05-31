package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/app"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/database"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/routes"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/session"
)

const PORT = 8080

func main() {
	log.Println("Welcome to subscription service")
	app := app.App{
		DB:       database.InitDB(),
		Session:  session.InitSession(),
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		WG:       &sync.WaitGroup{},
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", PORT),
		Handler: routes.InitRoutes(&app),
	}

	app.InfoLog.Printf("Starting web server on port %d\n", PORT)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
