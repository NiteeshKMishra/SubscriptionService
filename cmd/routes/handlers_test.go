package routes_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/database"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/routes"
)

var pageTests = []struct {
	name               string
	url                string
	expectedStatusCode int
	handler            http.HandlerFunc
	sessionData        map[string]any
	expectedHTML       string
}{
	{
		name:               "home",
		url:                "/",
		expectedStatusCode: http.StatusOK,
		handler: func(w http.ResponseWriter, r *http.Request) {
			routes.HomePage(&testApp, w, r)
		},
		expectedHTML: `<h1 class="mt-5 text-center">Welcome to SubscriptionService</h1>`,
	},
	{
		name:               "login page",
		url:                "/login",
		expectedStatusCode: http.StatusOK,
		handler: func(w http.ResponseWriter, r *http.Request) {
			routes.LoginPage(&testApp, w, r)
		},
		expectedHTML: `<h1 class="mt-5">Login</h1>`,
	},
	{
		name:               "register page",
		url:                "/register",
		expectedStatusCode: http.StatusOK,
		handler: func(w http.ResponseWriter, r *http.Request) {
			routes.RegisterPage(&testApp, w, r)
		},
		expectedHTML: `<h1 class="mt-5">Register</h1>`,
	},
	{
		name:               "logout",
		url:                "/logout",
		expectedStatusCode: http.StatusSeeOther,
		handler: func(w http.ResponseWriter, r *http.Request) {
			routes.LogoutHandler(&testApp, w, r)
		},
		sessionData: map[string]any{
			"user_id": "00000000-0000-0000-0000-000000000001",
			"user":    database.User{},
		},
	},
}

func TestPages(t *testing.T) {
	for _, e := range pageTests {
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", e.url, nil)

		ctx := getCtx(req)
		req = req.WithContext(ctx)

		if len(e.sessionData) > 0 {
			for key, value := range e.sessionData {
				testApp.Session.Put(ctx, key, value)
			}
		}

		e.handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s failed: expected %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if len(e.expectedHTML) > 0 {
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("%s failed: expected to find %s, but did not", e.name, e.expectedHTML)
			}
		}
	}
}

func TestLoginHandler(t *testing.T) {
	postedData := url.Values{
		"email":    {"admin@example.com"},
		"password": {"abc1234"},
	}

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routes.LoginHandler(&testApp, w, r)
	})
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Error("wrong code returned")
	}

	if !testApp.Session.Exists(ctx, "user_id") {
		t.Error("did not find userID in session")
	}
}

func TestSubscriptionHandler(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/subscribe?id=00000000-0000-0000-0000-000000000005", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	testApp.Session.Put(ctx, "user", database.User{
		ID:        "00000000-0000-0000-0000-000000000001",
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "User",
		Active:    true,
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routes.SubscriptionHandler(&testApp, w, r)
	})
	handler.ServeHTTP(rr, req)

	testApp.WG.Wait()

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status code of statusseeother, but got %d", rr.Code)
	}

	user := testApp.Session.Get(req.Context(), "user").(*database.User)
	if user == nil {
		t.Error("expected user in the session")
	}
	if user.Plan == nil {
		t.Errorf("expected plan to be added for user %s", user.Email)
	}
	if user.Plan.PlanName != "SILVER" {
		t.Errorf("expected plan to be SILVER BUT got %s", user.Plan.PlanName)
	}
}
