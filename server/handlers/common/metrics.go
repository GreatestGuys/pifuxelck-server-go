package common

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

// InstallHandler installs a handler on a given router. This method will also
// automatically instrument the handler with prometheus which will automatically
// report QPS, and quartiles for latency and response size.
func InstallHandler(r *mux.Router, path string, f http.HandlerFunc) *mux.Route {
	return r.HandleFunc(path, prometheus.InstrumentHandlerFunc(path, f))
}
