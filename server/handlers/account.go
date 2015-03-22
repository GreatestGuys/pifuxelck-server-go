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
	s := r.PathPrefix("/account/").Subrouter()
	s.HandleFunc("/login", accountLogin).Methods("POST")
	s.HandleFunc("/register", accountRegister).Methods("POST")
}

func accountLogin(w http.ResponseWriter, r *http.Request) {
	msg, err := common.RequestMessage(r)
	if err != nil {
		common.RespondClientError(w, &common.Errors{Meta: err})
		return
	}

	log.Debugf("Attempting to look up user %#v.", msg.User.DisplayName)
	id, userErr := models.UserLookupByPassword(*msg.User)
	if userErr != nil {
		common.RespondClientError(w, &common.Errors{User: userErr})
		return
	}

	log.Debugf("Creating new auth token for %#v.", msg.User.DisplayName)
	auth, metaErr := models.NewAuthToken(id)
	if metaErr != nil {
		common.RespondClientError(w, &common.Errors{Meta: metaErr})
		return
	}

	log.Infof("Successfully logged in as user %#v.", msg.User.DisplayName)
	common.RespondSuccess(w, &common.Message{
		User: &models.User{ID: id, DisplayName: msg.User.DisplayName},
		Meta: &models.Meta{Auth: auth},
	})
}

func accountRegister(w http.ResponseWriter, r *http.Request) {
	msg, err := common.RequestMessage(r)
	if err != nil {
		common.RespondClientError(w, &common.Errors{Meta: err})
		return
	}

	log.Debugf("Attempting to register new user %#v.", msg.User.DisplayName)
	user, userErr := models.CreateUser(*msg.User)
	if userErr != nil {
		log.Debugf("Failed to register user %#v.", user.DisplayName, user.ID)
		common.RespondClientError(w, &common.Errors{User: userErr})
		return
	}

	log.Infof("Successfully registered new user %#v (%v).", user.DisplayName, user.ID)
	common.RespondSuccess(w, &common.Message{User: user})
}
