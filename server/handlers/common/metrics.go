package common

import (
	"net/http"

	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

var metricUncaughtPanics = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "uncaught_panics",
		Help: "The number of uncaught panics per endpoint.",
	},
	[]string{"handler"},
)

func init() {
	prometheus.MustRegister(metricUncaughtPanics)
}

// InstallHandler installs a handler on a given router. This method will also
// automatically instrument the handler with prometheus which will automatically
// report QPS, and quartiles for latency and response size.
func InstallHandler(r *mux.Router, path string, f http.HandlerFunc) *mux.Route {
	metricUncaughtPanics.WithLabelValues(path).Add(0)

	wrapper := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				metricUncaughtPanics.WithLabelValues(path).Inc()
				log.Errorf("Uncaught panic, %v", r)
			}
		}()
		f(w, r)
	}

	return r.HandleFunc(path, prometheus.InstrumentHandlerFunc(path, wrapper))
}
