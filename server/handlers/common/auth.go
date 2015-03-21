package common

import (
	"fmt"
	"net/http"

	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
)

type authHandler func(string, http.ResponseWriter, *http.Request)

func (h authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.RawQuery != "password" {
		log.Warnf("Unauthenticated attempt to access restricted resource %v.", r.URL)

		w.WriteHeader(403)
		fmt.Fprintf(w, "You need to be authenticated!")
		return
	}

	userID := "myUserId"

	log.Verbosef("Successfully authenticated as user %v.", userID)
	h(userID, w, r)
}

// AuthHandlerFunc takes an function that takes a user ID, an
// http.ResponseWriter, and an http.Request and returns an http.Handler that
// will invoke the supplied function when a properly authenticated request is
// made, and returns a 403 error.
func AuthHandlerFunc(h func(string, http.ResponseWriter, *http.Request)) http.Handler {
	return authHandler(h)
}
