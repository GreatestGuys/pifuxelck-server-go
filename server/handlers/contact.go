package handlers

import (
	"net/http"

	"github.com/GreatestGuys/pifuxelck-server-go/server/handlers/common"
	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
	"github.com/GreatestGuys/pifuxelck-server-go/server/models"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricContactLookupSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "contact_lookup_success",
		Help: "The number of successful contact lookups.",
	})

	metricContactLookupFailure = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "contact_lookup_failure",
		Help: "The number of unsuccessful contact lookups.",
	})
)

func init() {
	prometheus.MustRegister(metricContactLookupFailure)
	prometheus.MustRegister(metricContactLookupSuccess)
}

// InstallContactHandlers takes a gorilla router and installs /contact/*
// endpoints on it.
func InstallContactHandlers(r *mux.Router) {
	common.InstallHandler(r, "/contacts/lookup/{displayName}", contactLookup).
		Methods("GET")
}

var contactLookup = common.AuthHandlerFunc(func(_ string, w http.ResponseWriter, r *http.Request) {
	displayName := mux.Vars(r)["displayName"]
	user, userErr := models.ContactLookup(displayName)

	if userErr != nil {
		metricContactLookupFailure.Inc()
		log.Debugf("Lookup for contact %#v failed.", displayName)
		common.RespondClientError(w, &models.Errors{User: userErr})
		return
	}

	metricContactLookupSuccess.Inc()
	log.Infof("Successful lookup of contact %#v.", displayName)
	common.RespondSuccess(w, &models.Message{User: user})
})
