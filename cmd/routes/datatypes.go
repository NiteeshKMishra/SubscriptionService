package routes

import (
	"time"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/database"
)

type TemplateData struct {
	Data          map[string]any
	Flash         string
	Warning       string
	Error         string
	Authenticated bool
	Now           time.Time
	User          *database.User
}
