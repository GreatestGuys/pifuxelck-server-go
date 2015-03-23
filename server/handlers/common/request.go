package common

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
	"github.com/GreatestGuys/pifuxelck-server-go/server/models"
)

// RequestMessage retrieves a Message from a given http request.
func RequestMessage(r *http.Request) (*models.Message, *models.Errors) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warnf("Unable to fully read request body due to error %v", err.Error())
		return nil, &models.Errors{App: []string{"Invalid request body."}}
	}

	msg := &models.Message{}
	err = json.Unmarshal(b, &msg)
	if err != nil {
		log.Warnf("Unable to unmarshal request body due to error %v", err.Error())
		return nil, &models.Errors{
			App: []string{"Unable to unmarshal request body."},
		}
	}

	return msg, nil
}

func RequestUserMessage(r *http.Request) (*models.User, *models.Errors) {
	msg, err := RequestMessage(r)
	if err != nil {
		return nil, err
	}

	if msg.User == nil {
		return nil, &models.Errors{User: &models.UserError{
			DisplayName: []string{"No user given."},
		}}
	}

	return msg.User, nil
}

func RequestNewGameMessage(r *http.Request) (*models.NewGame, *models.Errors) {
	msg, err := RequestMessage(r)
	if err != nil {
		return nil, err
	}

	if msg.NewGame == nil {
		return nil, &models.Errors{NewGame: &models.NewGameError{
			Label: []string{"No new game given."},
		}}
	}

	return msg.NewGame, nil
}
