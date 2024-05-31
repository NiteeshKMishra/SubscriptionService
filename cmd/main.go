package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/app"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/database"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/routes"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/session"
)

const PORT = 8080

func main() {
	log.Println("Welcome to subscription service")

	mutex := &sync.Mutex{}
	db := database.InitDB(mutex)
	app := app.App{
		DB:       db,
		Session:  session.InitSession(),
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		WG:       &sync.WaitGroup{},
		MU:       mutex,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", PORT),
		Handler: routes.InitRoutes(&app),
	}

	go listenForShutdown(&app)
	app.InfoLog.Printf("Starting web server on port %d\n", PORT)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func listenForShutdown(app *app.App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	shutDown(app)
	os.Exit(0)
}

func shutDown(app *app.App) {
	app.InfoLog.Println("Cleanup before shutdown, wait to background routines")
	app.WG.Wait()
	if app.DB != nil {
		app.DB.Close()
	}
	app.InfoLog.Println("Application shutdown")
}
