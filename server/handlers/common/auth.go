package common

import (
	"net/http"

	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
	"github.com/GreatestGuys/pifuxelck-server-go/server/models"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricAuthSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auth_success",
		Help: "The number of valid authenticated requests.",
	})

	metricAuthFailure = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auth_failure",
		Help: "The number of invalid authenticated requests.",
	})
)

func init() {
	prometheus.MustRegister(metricAuthFailure)
	prometheus.MustRegister(metricAuthSuccess)
}

// AuthHandlerFunc takes an function that takes a user ID, an
// http.ResponseWriter, and an http.Request and returns an http.Handler that
// will invoke the supplied function when a properly authenticated request is
// made, and returns a 403 error.
func AuthHandlerFunc(h func(int64, http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("x-pifuxelck-auth")
		userID, err := models.AuthTokenLookup(auth)
		if err != nil {
			metricAuthFailure.Inc()
			log.Debugf("Invalid authentication token %#v.", auth)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		metricAuthSuccess.Inc()
		log.Debugf("Successfully authenticated as user %v.", userID)
		h(userID, w, r)
	}
}
