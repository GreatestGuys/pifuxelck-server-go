package handlers

import (
	"net/http"

	"github.com/GreatestGuys/pifuxelck-server-go/server/handlers/common"
	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
	"github.com/GreatestGuys/pifuxelck-server-go/server/models"
	"github.com/gorilla/mux"
)

// Account creates an returns a router that is capable of handling all requests
// to the /account/* path.
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

	log.Debugf("Attempting to login as user %v.", msg.User.DisplayName)
	auth, userErr := models.UserLogin(*msg.User)
	if userErr != nil {
		common.RespondClientError(w, &common.Errors{User: userErr})
		return
	}

	log.Debugf("Successfully logged in as user %v.", msg.User.DisplayName)
	common.RespondSuccess(w, &common.Message{Meta: &common.Meta{Auth: auth}})
}

func accountRegister(w http.ResponseWriter, r *http.Request) {
	msg, err := common.RequestMessage(r)
	if err != nil {
		common.RespondClientError(w, &common.Errors{Meta: err})
		return
	}

	log.Debugf("Attempting to register new user %v.", msg.User.DisplayName)
	user, userErr := models.CreateUser(*msg.User)
	if userErr != nil {
		common.RespondClientError(w, &common.Errors{User: userErr})
		return
	}

	log.Debugf("Successfully registered new user %v (%v).", user.DisplayName, user.ID)
	common.RespondSuccess(w, &common.Message{User: user})
}
