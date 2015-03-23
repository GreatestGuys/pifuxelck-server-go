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
	metricGameCreateSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "game_create_success",
		Help: "The number of created games.",
	})

	metricGameCreateFailure = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "game_create_failure",
		Help: "The number of failed attempts to create a game.",
	})
)

func init() {
	prometheus.MustRegister(metricGameCreateFailure)
	prometheus.MustRegister(metricGameCreateSuccess)
}

// InstallGameHandlers takes a gorilla router and installs /games/* endpoints on
// it.
func InstallGameHandlers(r *mux.Router) {
	common.InstallHandler(r, "/games/new", gameCreate).Methods("POST")
}

var gameCreate = common.AuthHandlerFunc(func(id string, w http.ResponseWriter, r *http.Request) {
	newGame, err := common.RequestNewGameMessage(r)
	if err != nil {
		metricGameCreateFailure.Inc()
		common.RespondClientError(w, err)
		return
	}

	log.Debugf("Attempting to start new game.")
	errors := models.CreateGame(id, *newGame)
	if errors != nil {
		metricGameCreateFailure.Inc()
		log.Debugf("Failed to create new game.")
		common.RespondClientError(w, errors)
		return
	}

	metricGameCreateSuccess.Inc()
	log.Infof("User %v created new game.", id)
	common.RespondSuccessNoContent(w)
})
