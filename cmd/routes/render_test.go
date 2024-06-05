package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/routes"
)

func TestAddDefaultData(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	testApp.Session.Put(ctx, "flash", "flash")
	testApp.Session.Put(ctx, "warning", "warning")
	testApp.Session.Put(ctx, "error", "error")

	td := routes.AddDefaultTemplateData(&testApp, req, &routes.TemplateData{})

	if td.Flash != "flash" {
		t.Error("failed to get flash data")
	}

	if td.Warning != "warning" {
		t.Error("failed to get warning data")
	}

	if td.Error != "error" {
		t.Error("failed to get error data")
	}
}

func TestIsAuthenticated(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	auth := routes.IsAuthenticated(&testApp, req)

	if auth {
		t.Error("returns true for authenticated, when it should be false")
	}

	testApp.Session.Put(ctx, "user_id", "00000000-0000-0000-0000-000000000001")

	auth = routes.IsAuthenticated(&testApp, req)

	if !auth {
		t.Error("returns false for authenticated, when it should be true")
	}
}

func TestRender(t *testing.T) {
	rr := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	routes.Render(&testApp, rr, req, "home.page.gohtml", &routes.TemplateData{})

	if rr.Code != 200 {
		t.Error("failed to render page")
	}
}
