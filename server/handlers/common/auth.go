package common

import (
	"net/http"

	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
)

// AuthHandlerFunc takes an function that takes a user ID, an
// http.ResponseWriter, and an http.Request and returns an http.Handler that
// will invoke the supplied function when a properly authenticated request is
// made, and returns a 403 error.
func AuthHandlerFunc(h func(string, http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "password" {
			log.Warnf("Unauthenticated attempt to access restricted resource %v.", r.URL)

			w.WriteHeader(403)
			return
		}

		userID := "myUserId"

		log.Verbosef("Successfully authenticated as user %v.", userID)
		h(userID, w, r)
	}
}
