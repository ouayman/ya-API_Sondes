package monitoringService

import (
	"net/http"

	"../../Configuration"
	"../../Helper/Http"
)

type handler struct {
}

func newHandler() *handler {
	return &handler{}
}

func (obj handler) getInfo(w http.ResponseWriter, r *http.Request) error {
	return helperhttp.Respond(w, configuration.Get())
}
