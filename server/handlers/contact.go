package handlers

import (
	"net/http"

	"github.com/GreatestGuys/pifuxelck-server-go/server/handlers/common"
	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
	"github.com/GreatestGuys/pifuxelck-server-go/server/models"
	"github.com/gorilla/mux"
)

// InstallContactHandlers takes a gorilla router and installs /contact/*
// endpoints on it.
func InstallContactHandlers(r *mux.Router) {
	common.InstallHandler(r, "/contacts/lookup/{displayName}", contactLookup).
		Methods("GET")
}

var contactLookup = common.AuthHandlerFunc(func(_ int64, w http.ResponseWriter, r *http.Request) {
	displayName := mux.Vars(r)["displayName"]
	user, userErr := models.ContactLookup(displayName)

	if userErr != nil {
		log.Debugf("Lookup for contact %#v failed.", displayName)
		common.RespondClientError(w, &models.Errors{User: userErr})
		return
	}

	log.Infof("Successful lookup of contact %#v.", displayName)
	common.RespondSuccess(w, &models.Message{User: user})
})
