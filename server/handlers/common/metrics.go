package common

import (
	"net/http"
	"strconv"

	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricUncaughtPanics = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "uncaught_panics",
			Help: "The number of uncaught panics per endpoint.",
		},
		[]string{"handler"})

	metricEndpointQueries = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "endpoint_queries",
			Help: "The number of responses per endpoint per result code.",
		},
		[]string{"handler", "status"})
)

func init() {
	prometheus.MustRegister(metricUncaughtPanics)
	prometheus.MustRegister(metricEndpointQueries)
}

type responseWriterWrapper struct {
	handler string
	inner   http.ResponseWriter
}

func (w responseWriterWrapper) Header() http.Header {
	return w.inner.Header()
}

func (w responseWriterWrapper) Write(b []byte) (int, error) {
	return w.inner.Write(b)
}

func (w responseWriterWrapper) WriteHeader(status int) {
	metricEndpointQueries.WithLabelValues(w.handler, strconv.Itoa(status)).Inc()
	w.inner.WriteHeader(status)
}

// InstallHandler installs a handler on a given router. This method will also
// automatically instrument the handler with prometheus which will automatically
// report QPS, and quartiles for latency and response size.
func InstallHandler(r *mux.Router, path string, f http.HandlerFunc) *mux.Route {
	metricUncaughtPanics.WithLabelValues(path).Add(0)

	metricEndpointQueries.WithLabelValues(path, "200").Add(0)
	metricEndpointQueries.WithLabelValues(path, "204").Add(0)
	metricEndpointQueries.WithLabelValues(path, "403").Add(0)
	metricEndpointQueries.WithLabelValues(path, "422").Add(0)
	metricEndpointQueries.WithLabelValues(path, "500").Add(0)

	addCorsHeaders := func(w http.ResponseWriter) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "x-pifuxelck-auth")
	}

	wrapper := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				metricUncaughtPanics.WithLabelValues(path).Inc()
				log.Errorf("Uncaught panic, %v", r)
			}
		}()

		addCorsHeaders(w)
		f(responseWriterWrapper{handler: path, inner: w}, r)
	}

	cors := prometheus.InstrumentHandlerFunc(path,
		func(w http.ResponseWriter, r *http.Request) {
			wrappedWriter := responseWriterWrapper{handler: path, inner: w}
			addCorsHeaders(wrappedWriter)
			wrappedWriter.WriteHeader(http.StatusNoContent)
		})

	r.HandleFunc(path, cors).Methods("OPTIONS")
	return r.HandleFunc(path, prometheus.InstrumentHandlerFunc(path, wrapper))
}
