package handlers

import (
	"net/http"

	"github.com/GreatestGuys/pifuxelck-server-go/server/handlers/common"
	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
	"github.com/GreatestGuys/pifuxelck-server-go/server/models"
	"github.com/gorilla/mux"
)

// InstallAccountHandlers takes a gorilla router and installs /account/*
// endpoints on it.
func InstallAccountHandlers(r *mux.Router) {
	r.HandleFunc("/account/login", accountLogin).Methods("POST")
	r.HandleFunc("/account/register", accountRegister).Methods("POST")
	r.HandleFunc("/account", accountUpdate).Methods("PUT")
}

func accountLogin(w http.ResponseWriter, r *http.Request) {
	msg, err := common.RequestMessage(r)
	if err != nil {
		common.RespondClientError(w, &models.Errors{Meta: err})
		return
	}

	log.Debugf("Attempting to look up user %#v.", msg.User.DisplayName)
	id, userErr := models.UserLookupByPassword(*msg.User)
	if userErr != nil {
		common.RespondClientError(w, &models.Errors{User: userErr})
		return
	}

	log.Debugf("Creating new auth token for %#v.", msg.User.DisplayName)
	auth, metaErr := models.NewAuthToken(id)
	if metaErr != nil {
		common.RespondClientError(w, &models.Errors{Meta: metaErr})
		return
	}

	log.Infof("Successfully logged in as user %#v.", msg.User.DisplayName)
	common.RespondSuccess(w, &models.Message{
		User: &models.User{ID: id, DisplayName: msg.User.DisplayName},
		Meta: &models.Meta{Auth: auth},
	})
}

func accountRegister(w http.ResponseWriter, r *http.Request) {
	msg, err := common.RequestMessage(r)
	if err != nil {
		common.RespondClientError(w, &models.Errors{Meta: err})
		return
	}

	log.Debugf("Attempting to register new user %#v.", msg.User.DisplayName)
	user, userErr := models.CreateUser(*msg.User)
	if userErr != nil {
		log.Debugf("Failed to register user %#v.", user.DisplayName, user.ID)
		common.RespondClientError(w, &models.Errors{User: userErr})
		return
	}

	log.Infof("Successfully registered new user %#v (%v).", user.DisplayName, user.ID)
	common.RespondSuccess(w, &models.Message{User: user})
}

var accountUpdate = common.AuthHandlerFunc(func(id string, w http.ResponseWriter, r *http.Request) {
	msg, err := common.RequestMessage(r)
	if err != nil {
		common.RespondClientError(w, &models.Errors{Meta: err})
		return
	}

	log.Debugf("Attempting to update password for %#v.", msg.User.DisplayName)

	// Override any ID given in the JSON request body with the actual
	// authenticated user ID.
	msg.User.ID = id
	user, userErr := models.UserSetPassword(*msg.User)
	if userErr != nil {
		log.Debugf("Failed to update password, %v.", userErr.Error())
		common.RespondClientError(w, &models.Errors{User: userErr})
		return
	}

	log.Infof("Successfully updated password of %#v.", user.DisplayName)
	common.RespondSuccess(w, &models.Message{User: user})
})
