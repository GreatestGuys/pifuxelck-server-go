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

	metricGameHistorySuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "game_history_success",
		Help: "The number of successful history queries.",
	})

	metricGameHistoryFailure = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "game_history_failure",
		Help: "The number of failed attempts to query history.",
	})

	metricGameCompletedAtFailure = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "game_completed_at_update_failed",
		Help: "The number of times the query to update a game's finish time fails.",
	})
)

func init() {
	prometheus.MustRegister(metricGameCreateFailure)
	prometheus.MustRegister(metricGameCreateSuccess)
	prometheus.MustRegister(metricGameInboxFailure)
	prometheus.MustRegister(metricGameInboxSuccess)
	prometheus.MustRegister(metricGamePlayFailure)
	prometheus.MustRegister(metricGamePlaySuccess)
	prometheus.MustRegister(metricGameHistoryFailure)
	prometheus.MustRegister(metricGameHistorySuccess)
	prometheus.MustRegister(metricGameCompletedAtFailure)
}

// InstallGameHandlers takes a gorilla router and installs /games/* endpoints on
// it.
func InstallGameHandlers(r *mux.Router) {
	common.InstallHandler(r, "/games/new", gameCreate).Methods("POST")
	common.InstallHandler(r, "/games/inbox", gameInbox).Methods("GET")
	common.InstallHandler(r, "/games/play/{id:[0-9]+}", gamePlay).Methods("PUT")
	common.InstallHandler(r, "/games/since/{id:[0-9]+}", gameHistory).
		Methods("GET")
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
	models.ReapExpiredTurns()

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

	// Check if the game needs to have it's completed at time updated.
	errors = models.UpdateGameCompletedAtTime(gameID)
	if errors != nil {
		// This should not be returned to the user, as their turn has already been
		// successfully added, instead just log this as an error, update the
		// corresponding metric and hope that someone checks the logs.
		log.Errorf("Unable to update game completed at time!")
		metricGameCompletedAtFailure.Inc()
	}

	metricGamePlaySuccess.Inc()
	log.Infof("User %v took their turn in game %v.", userID, gameID)
	common.RespondSuccessNoContent(w)
})

var gameHistory = common.AuthHandlerFunc(func(userID string, w http.ResponseWriter, r *http.Request) {
	sinceID := mux.Vars(r)["id"]

	log.Debugf("User %v is requesting history since %v.", userID, sinceID)

	games, errors := models.CompletedGames(userID, sinceID)
	if errors != nil {
		metricGameHistoryFailure.Inc()
		common.RespondClientError(w, errors)
		return
	}

	metricGameHistorySuccess.Inc()
	log.Infof("User %v looked up history since %v.", userID, sinceID)
	common.RespondSuccess(w, &models.Message{Games: games})
})
