package app

import (
	"database/sql"
	"log"
	"sync"

	"github.com/alexedwards/scs/v2"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/database"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/emailer"
)

type App struct {
	DB            *sql.DB
	Session       *scs.SessionManager
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	WG            *sync.WaitGroup
	MU            *sync.Mutex
	Models        database.Models
	Mailer        emailer.Mail
	ErrorChan     chan error
	ErrorChanDone chan bool
}
