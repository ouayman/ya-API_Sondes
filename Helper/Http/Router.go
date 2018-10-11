package helperhttp

import (
	"net/http"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.Handler
}

type Routes []Route
