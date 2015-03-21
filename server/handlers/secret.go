package handlers

import (
	"fmt"
	"net/http"

	"github.com/GreatestGuys/pifuxelck-server-go/server/handlers/common"
	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
)

var Secret = common.AuthHandlerFunc(func (id string, w http.ResponseWriter, r *http.Request) {
	log.Verbosef("Accessing secret content as user %v.", id)
	fmt.Fprintf(w, "Welcome to the secret page, %v!", id)
})
