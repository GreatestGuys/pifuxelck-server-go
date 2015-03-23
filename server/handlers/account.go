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
	metricLoginSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "account_login_success",
		Help: "The number of successful logins.",
	})

	metricLoginFailure = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "account_login_failure",
		Help: "The number of unsuccessful logins.",
	})

	metricRegisterSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "account_register_success",
		Help: "The number of successful account registrations.",
	})

	metricRegisterFailure = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "account_register_failure",
		Help: "The number of unsuccessful account registrations.",
	})

	metrictAccountUpdateSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "account_update_success",
		Help: "The number of successful account update requests.",
	})

	metrictAccountUpdateFailure = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "account_update_failure",
		Help: "The number of unsuccessful account update requests.",
	})
)

func init() {
	prometheus.MustRegister(metricLoginFailure)
	prometheus.MustRegister(metricLoginSuccess)
	prometheus.MustRegister(metricRegisterFailure)
	prometheus.MustRegister(metricRegisterSuccess)
	prometheus.MustRegister(metrictAccountUpdateFailure)
	prometheus.MustRegister(metrictAccountUpdateSuccess)
}

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
		metricLoginFailure.Inc()
		common.RespondClientError(w, err)
		return
	}

	log.Debugf("Attempting to look up user %#v.", user.DisplayName)
	id, userErr := models.UserLookupByPassword(*user)
	if userErr != nil {
		metricLoginFailure.Inc()
		common.RespondClientError(w, &models.Errors{User: userErr})
		return
	}

	log.Debugf("Creating new auth token for %#v.", user.DisplayName)
	auth, errors := models.NewAuthToken(id)
	if errors != nil {
		metricLoginFailure.Inc()
		common.RespondClientError(w, errors)
		return
	}

	metricLoginSuccess.Inc()
	log.Infof("Successfully logged in as user %#v.", user.DisplayName)
	common.RespondSuccess(w, &models.Message{
		User: &models.User{ID: id, DisplayName: user.DisplayName},
		Meta: &models.Meta{Auth: auth},
	})
}

func accountRegister(w http.ResponseWriter, r *http.Request) {
	user, err := common.RequestUserMessage(r)
	if err != nil {
		metricRegisterFailure.Inc()
		common.RespondClientError(w, err)
		return
	}

	log.Debugf("Attempting to register new user %#v.", user.DisplayName)
	user, userErr := models.CreateUser(*user)
	if userErr != nil {
		metricRegisterFailure.Inc()
		log.Debugf("Failed to register user %#v.", user.DisplayName, user.ID)
		common.RespondClientError(w, &models.Errors{User: userErr})
		return
	}

	metricRegisterSuccess.Inc()
	log.Infof("Successfully registered new user %#v (%v).", user.DisplayName, user.ID)
	common.RespondSuccess(w, &models.Message{User: user})
}

var accountUpdate = common.AuthHandlerFunc(func(id string, w http.ResponseWriter, r *http.Request) {
	user, err := common.RequestUserMessage(r)
	if err != nil {
		metrictAccountUpdateFailure.Inc()
		common.RespondClientError(w, err)
		return
	}

	log.Debugf("Attempting to update password for %#v.", user.DisplayName)

	// Override any ID given in the JSON request body with the actual
	// authenticated user ID.
	user.ID = id
	user, userErr := models.UserSetPassword(*user)
	if userErr != nil {
		metrictAccountUpdateFailure.Inc()
		log.Debugf("Failed to update password, %v.", userErr.Error())
		common.RespondClientError(w, &models.Errors{User: userErr})
		return
	}

	metrictAccountUpdateSuccess.Inc()
	log.Infof("Successfully updated password of %#v.", user.DisplayName)
	common.RespondSuccess(w, &models.Message{User: user})
})
