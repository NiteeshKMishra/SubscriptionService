package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/emailer"
)

func (app *App) ListenForMail() {
	for {
		select {
		case msg := <-app.Mailer.MailerChan:
			go app.Mailer.SendMail(msg, app.Mailer.ErrorChan)
		case err := <-app.Mailer.ErrorChan:
			app.ErrorLog.Printf("error sending email: %s", err.Error())
		case <-app.Mailer.DoneChan:
			return
		}
	}
}

func (app *App) SendEmail(msg emailer.Message) {
	app.WG.Add(1)
	app.Mailer.MailerChan <- msg
}

func (app *App) ListenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutDown()
	os.Exit(0)
}

func (app *App) shutDown() {
	app.InfoLog.Println("Cleanup before shutdown, wait to background routines")
	app.WG.Wait()
	if app.DB != nil {
		app.DB.Close()
	}

	app.Mailer.DoneChan <- true
	close(app.Mailer.MailerChan)
	close(app.Mailer.ErrorChan)
	close(app.Mailer.DoneChan)

	app.InfoLog.Println("Application shutdown")
}
