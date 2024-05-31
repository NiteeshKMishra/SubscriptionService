package app

import (
	"database/sql"
	"log"
	"sync"

	"github.com/alexedwards/scs/v2"
)

type App struct {
	DB       *sql.DB
	Session  *scs.SessionManager
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	WG       *sync.WaitGroup
	MU       *sync.Mutex
}
