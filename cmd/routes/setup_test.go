package routes_test

import (
	"context"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/app"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/constants"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/database"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/emailer"
)

var testApp app.App

func TestMain(m *testing.M) {
	gob.Register(database.User{})

	// set up session
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	testApp = app.App{
		Session:       session,
		DB:            nil,
		InfoLog:       log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog:      log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		Models:        database.TestNew(nil),
		WG:            &sync.WaitGroup{},
		ErrorChan:     make(chan error),
		ErrorChanDone: make(chan bool),
	}

	// create a dummy mailer
	errorChan := make(chan error)
	mailerChan := make(chan emailer.Message, 100)
	mailerDoneChan := make(chan bool)

	testApp.Mailer = emailer.Mail{
		Wait:       testApp.WG,
		ErrorChan:  errorChan,
		MailerChan: mailerChan,
		DoneChan:   mailerDoneChan,
	}

	go func() {
		for {
			select {
			case <-testApp.Mailer.MailerChan:
				testApp.WG.Done()
			case <-testApp.Mailer.ErrorChan:
			case <-testApp.Mailer.DoneChan:
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case err := <-testApp.ErrorChan:
				testApp.ErrorLog.Println(err)
			case <-testApp.ErrorChanDone:
				return
			}
		}
	}()

	existingTempPath := constants.PathToTemplates
	existingAssetPath := constants.PathToAssets

	constants.PathToTemplates = "../templates"
	constants.PathToAssets = "../../assets"

	os.Exit(m.Run())

	constants.PathToTemplates = existingTempPath
	constants.PathToAssets = existingAssetPath
}

func getCtx(req *http.Request) context.Context {
	ctx, err := testApp.Session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
