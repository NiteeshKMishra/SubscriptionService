package routes_test

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/routes"
)

var validRoutes = []string{
	"/",
	"/login",
	"/logout",
	"/register",
	"/activate",
	"/user/plans",
	"/user/subscribe",
}

var inValidRoutes = []string{
	"/todo",
	"/user/foo",
}

func TestRoutes(t *testing.T) {
	testRoutes := routes.InitRoutes(&testApp)

	chiRoutes := testRoutes.(chi.Router)

	for _, route := range validRoutes {
		routeExists(t, chiRoutes, route)
	}

	for _, route := range inValidRoutes {
		routeNotExists(t, chiRoutes, route)
	}
}

func routeExists(t *testing.T, routes chi.Router, route string) {
	found := false

	_ = chi.Walk(routes, func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if route == foundRoute {
			found = true
		}
		return nil
	})

	if !found {
		t.Errorf("did not find %s in registered routes", route)
	}
}

func routeNotExists(t *testing.T, routes chi.Router, route string) {
	found := false

	_ = chi.Walk(routes, func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if route == foundRoute {
			found = true
		}
		return nil
	})

	if found {
		t.Errorf("find %s in registered routes, should have not exists", route)
	}
}
