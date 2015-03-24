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

// RequestUserMessage extracts and returns a User model from the request body
// and returns an error if unable to do so.
func RequestUserMessage(r *http.Request) (*models.User, *models.Errors) {
	msg, err := RequestMessage(r)
	if err != nil {
		return nil, err
	}

	if msg.User == nil {
		return nil, &models.Errors{App: []string{"No user object in request body."}}
	}

	return msg.User, nil
}

// RequestNewGameMessage extracts and returns a NewGame model from the request
// body and returns an error if unable to do so.
func RequestNewGameMessage(r *http.Request) (*models.NewGame, *models.Errors) {
	msg, err := RequestMessage(r)
	if err != nil {
		return nil, err
	}

	if msg.NewGame == nil {
		return nil, &models.Errors{
			App: []string{"No new_game object in request body."}}
	}

	return msg.NewGame, nil
}

// RequestTurnMessage extracts and returns a Turn model from the request body
// and returns an error if unable to do so.
func RequestTurnMessage(r *http.Request) (*models.Turn, *models.Errors) {
	msg, err := RequestMessage(r)
	if err != nil {
		return nil, err
	}

	if msg.Turn == nil {
		return nil, &models.Errors{App: []string{"No turn object in request body."}}
	}

	return msg.Turn, nil
}
