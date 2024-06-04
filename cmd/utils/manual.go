package utils

import (
	"fmt"
	"os"
	"path"

	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/app"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/constants"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/database"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/emailer"
)

func SendManual(app *app.App, user database.User, plan *database.Plan) {
	defer app.WG.Done()
	manualPath := fmt.Sprintf("%s/tmp/%s_manual.pdf", constants.PathToAssets, user.ID)
	dirPath := path.Join(constants.PathToAssets, "tmp")

	// makes sure path to manual.pdf exists
	_, err := os.Stat(dirPath)
	if err != nil {
		if !os.IsNotExist(err) {
			app.ErrorLog.Printf("unable to generate manual: %s", err.Error())
			app.ErrorChan <- err
			return
		} else {
			err = os.Mkdir(dirPath, constants.DirPermission)
			if err != nil {
				app.ErrorLog.Printf("unable to generate manual: %s", err.Error())
				app.ErrorChan <- err
				return
			}
		}
	}

	pdf := generateManual(user, plan)
	err = pdf.OutputFileAndClose(manualPath)
	if err != nil {
		app.ErrorLog.Printf("unable to generate manual: %s", err.Error())
		app.ErrorChan <- err
		return
	}

	msg := emailer.Message{
		To:      user.Email,
		Subject: "Your manual",
		Data:    "Your user manual is attached",
		AttachmentMap: map[string]string{
			"Manual.pdf": manualPath,
		},
	}

	app.SendEmail(msg)
}

func generateManual(user database.User, plan *database.Plan) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetMargins(10, 13, 10)

	importer := gofpdi.NewImporter()

	t := importer.ImportPage(pdf, fmt.Sprintf("%s/manual.pdf", constants.PathToAssets), 1, "/MediaBox")
	pdf.AddPage()

	importer.UseImportedTemplate(pdf, t, 0, 0, 215.9, 0)

	pdf.SetX(75)
	pdf.SetY(150)

	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 4, fmt.Sprintf("%s %s", user.FirstName, user.LastName), "", "C", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 4, fmt.Sprintf("%s User Guide", plan.PlanName), "", "C", false)

	return pdf
}
