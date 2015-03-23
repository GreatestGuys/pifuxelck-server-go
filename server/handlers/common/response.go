package common

import (
	"encoding/json"
	"net/http"

	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
	"github.com/GreatestGuys/pifuxelck-server-go/server/models"
)

func setContentTypeToJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
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
func RespondClientError(w http.ResponseWriter, r *models.Errors) {
	b, err := json.Marshal(models.Message{Errors: r})
	if err != nil {
		log.Errorf("Unable to marshal response %v, due to error %v.", r, err.Error())
		RespondServerError(w)
		return
	}

	setContentTypeToJson(w)
	w.WriteHeader(422)
	w.Write(b)
}

// RespondSuccess signals to the client that the query was a success and returns
// an encoded Message.
func RespondSuccess(w http.ResponseWriter, r *models.Message) {
	b, err := json.Marshal(r)
	if err != nil {
		log.Errorf("Unable to marshal response %v, due to error %v.", r, err.Error())
		RespondServerError(w)
		return
	}

	setContentTypeToJson(w)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// RespondSuccessNoContent responds with an empty success message. If there is
// no information to communicate to the client then this method is preferred
// over RespondSuccess since it saves bandwidth and is marginally faster in
// building the actual response.
func RespondSuccessNoContent(w http.ResponseWriter) {
	setContentTypeToJson(w)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
