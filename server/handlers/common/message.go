package common

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
	"github.com/GreatestGuys/pifuxelck-server-go/server/models"
	"github.com/GreatestGuys/pifuxelck-server-go/server/models/common"
)

// Message corresponds to the top level JSON object that is returned by all
// end points.
type Message struct {
	Errors *Errors      `json:"errors,omitempty"`
	User   *models.User `json:"user,omitempty"`
	Meta   *Meta        `json:"meta,omitempty"`
}

type Errors struct {
	User *models.UserError `json:"user,omitempty"`
	Meta *MetaError        `json:"meta,omitempty"`
}

// Meta encodes meta data that does not correspond to any particular model.
type Meta struct {
	Auth string `json:"auth,omitempty"`
}

// MetaError encodes errors that do not correspond to any particular model.
type MetaError struct {
	Encoding []string `json:"encoding,omitempty"`
}

func (e MetaError) Error() string {
	return common.ModelErrorHelper(e)
}

// RespondServerError signals to the client that the server encountered an
// unexpected error and was unable to respond. This method should only be used
// when the server encounters an unrecoverable error. To signal that error lies
// in the request from the client, RespondClientError should be used instead.
func RespondServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

// RespondSuccess signals to the client that the query was a success and returns
// an encoded Message.
func RespondClientError(w http.ResponseWriter, r *Errors) {
	b, err := json.Marshal(Message{Errors: r})
	if err != nil {
		log.Errorf("Unable to marshal response %v, due to error %v.", r, err.Error())
		RespondServerError(w)
		return
	}

	w.WriteHeader(422)
	w.Write(b)
}

// RespondSuccess signals to the client that the query was a success and returns
// an encoded Message.
func RespondSuccess(w http.ResponseWriter, r *Message) {
	b, err := json.Marshal(r)
	if err != nil {
		log.Errorf("Unable to marshal response %v, due to error %v.", r, err.Error())
		RespondServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// RespondSuccessNoContent responds with an empty success message. If there is
// no information to communicate to the client then this method is preferred
// over RespondSuccess since it saves bandwidth and is marginally faster in
// building the actual response.
func RespondSuccessNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

// RequestMessage retrieves a Message from a given http request.
func RequestMessage(r *http.Request) (*Message, *MetaError) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warnf("Unable to fully read request body due to error %v", err.Error())
		return nil, &MetaError{Encoding: []string{"Invalid request body."}}
	}

	msg := &Message{}
	err = json.Unmarshal(b, &msg)
	if err != nil {
		log.Warnf("Unable to unmarshal request body due to error %v", err.Error())
		return nil, &MetaError{Encoding: []string{"Unable to unmarshal request body."}}
	}

	return msg, nil
}
