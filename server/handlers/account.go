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
	common.InstallHandler(r, "/account/login", accountLogin).Methods("POST")
	common.InstallHandler(r, "/account/register", accountRegister).Methods("POST")
	common.InstallHandler(r, "/account", accountUpdate).Methods("PUT")
}

func accountLogin(w http.ResponseWriter, r *http.Request) {
	user, err := common.RequestUserMessage(r)
	if err != nil {
		common.RespondClientError(w, err)
		return
	}

	log.Debugf("Attempting to look up user %#v.", user.DisplayName)
	id, userErr := models.UserLookupByPassword(*user)
	if userErr != nil {
		common.RespondClientError(w, &models.Errors{User: userErr})
		return
	}

	log.Debugf("Creating new auth token for %#v.", user.DisplayName)
	auth, errors := models.NewAuthToken(id)
	if errors != nil {
		common.RespondClientError(w, errors)
		return
	}

	log.Infof("Successfully logged in as user %#v.", user.DisplayName)
	common.RespondSuccess(w, &models.Message{
		User: &models.User{ID: id, DisplayName: user.DisplayName},
		Meta: &models.Meta{Auth: auth},
	})
}

func accountRegister(w http.ResponseWriter, r *http.Request) {
	origUser, err := common.RequestUserMessage(r)
	if err != nil {
		common.RespondClientError(w, err)
		return
	}

	log.Debugf("Attempting to register new user %#v.", origUser.DisplayName)
	user, userErr := models.CreateUser(*origUser)
	if userErr != nil {
		log.Debugf("Failed to register user %#v.", origUser.DisplayName)
		common.RespondClientError(w, &models.Errors{User: userErr})
		return
	}

	log.Debugf("Creating new auth token for %#v.", user.DisplayName)
	auth, errors := models.NewAuthToken(user.ID)
	if errors != nil {
		common.RespondClientError(w, errors)
		return
	}

	log.Infof("Successfully registered new user %#v (%v).", user.DisplayName, user.ID)
	common.RespondSuccess(w, &models.Message{
		User: user,
		Meta: &models.Meta{Auth: auth},
	})
}

var accountUpdate = common.AuthHandlerFunc(func(id int64, w http.ResponseWriter, r *http.Request) {
	user, err := common.RequestUserMessage(r)
	if err != nil {
		common.RespondClientError(w, err)
		return
	}

	log.Debugf("Attempting to update password for %#v.", user.DisplayName)

	// Override any ID given in the JSON request body with the actual
	// authenticated user ID.
	user.ID = id
	user, userErr := models.UserSetPassword(*user)
	if userErr != nil {
		log.Debugf("Failed to update password, %v.", userErr.Error())
		common.RespondClientError(w, &models.Errors{User: userErr})
		return
	}

	log.Infof("Successfully updated password of %#v.", user.DisplayName)
	common.RespondSuccess(w, &models.Message{User: user})
})
