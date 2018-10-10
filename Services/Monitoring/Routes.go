package monitoringService

import (
	"../../Helper/Http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func routes() helperhttp.Routes {
	handler := newHandler()

	return helperhttp.Routes{
		helperhttp.Route{
			Name:    "GetInfo",
			Method:  "GET",
			Pattern: "/info",
			Handler: helperhttp.ErrorFnHandler(handler.getInfo),
		},
		helperhttp.Route{
			Name:    "PrometheusHandler",
			Method:  "GET",
			Pattern: "/prometheus",
			Handler: promhttp.Handler(),
		},
	}
}
