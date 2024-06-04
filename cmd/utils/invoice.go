package utils

import (
	"errors"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/app"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/database"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/emailer"
)

func SendInvoice(app *app.App, user database.User, plan *database.Plan) {
	defer app.WG.Done()

	invoice := plan.PlanAmountFormatted
	if invoice == "" {
		app.ErrorChan <- errors.New("no price is added to plan")
		return
	}

	msg := emailer.Message{
		To:       user.Email,
		Subject:  "Your invoice",
		Data:     invoice,
		Template: "invoice",
	}

	app.SendEmail(msg)
}
