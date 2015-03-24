package common

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
	"github.com/GreatestGuys/pifuxelck-server-go/server/models"
)

// RequestMessage retrieves a Message from a given http request.
func RequestMessage(r *http.Request) (*models.Message, *models.MetaError) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warnf("Unable to fully read request body due to error %v", err.Error())
		return nil, &models.MetaError{Encoding: []string{"Invalid request body."}}
	}

	msg := &models.Message{}
	err = json.Unmarshal(b, &msg)
	if err != nil {
		log.Warnf("Unable to unmarshal request body due to error %v", err.Error())
		return nil, &models.MetaError{
			Encoding: []string{"Unable to unmarshal request body."},
		}
	}

	return msg, nil
}
