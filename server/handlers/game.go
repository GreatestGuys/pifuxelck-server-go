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

	metricGameInboxSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "game_inbox_success",
		Help: "The number of successful inbox queries.",
	})

	metricGameInboxFailure = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "game_inbox_failure",
		Help: "The number of failed attempts to retrieve the player's inbox.",
	})

	metricGamePlaySuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "game_play_success",
		Help: "The number of successful turns taken.",
	})

	metricGamePlayFailure = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "game_play_failure",
		Help: "The number of failed attempts to take a turn.",
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
	common.InstallHandler(r, "/games/inbox", gameInbox).Methods("GET")
	common.InstallHandler(r, "/games/play/{id:[0-9]+}", gamePlay).Methods("PUT")
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

var gameInbox = common.AuthHandlerFunc(func(id string, w http.ResponseWriter, r *http.Request) {
	log.Debugf("Attempting to query users inbox.")
	entries, errors := models.GetInboxEntriesForUser(id)
	if errors != nil {
		metricGameInboxFailure.Inc()
		log.Debugf("Failed to create new game.")
		common.RespondClientError(w, errors)
		return
	}

	metricGameInboxSuccess.Inc()
	log.Infof("User %v retrieved inbox.", id)
	common.RespondSuccess(w, &models.Message{InboxEntries: entries})
})

var gamePlay = common.AuthHandlerFunc(func(userID string, w http.ResponseWriter, r *http.Request) {
	turn, err := common.RequestTurnMessage(r)
	if err != nil {
		metricGamePlayFailure.Inc()
		common.RespondClientError(w, err)
		return
	}

	gameID := mux.Vars(r)["id"]

	log.Debugf("User %v is taking their turn in game %v.", userID, gameID)

	var takeTurn = func() *models.Errors {
		if turn.IsDrawing {
			return models.UpdateDrawingTurn(userID, gameID, turn.Drawing)
		} else {
			return models.UpdateLabelTurn(userID, gameID, turn.Label)
		}
	}

	errors := takeTurn()
	if errors != nil {
		metricGamePlayFailure.Inc()
		log.Debugf("User %v failed to take turn in game %v.", userID, gameID)
		common.RespondClientError(w, errors)
		return
	}

	metricGamePlaySuccess.Inc()
	log.Infof("User %v took their turn in game %v.", userID, gameID)
	common.RespondSuccessNoContent(w)
})
