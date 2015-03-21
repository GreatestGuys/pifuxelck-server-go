package server

import (
	"net/http"
	"strconv"

	"github.com/GreatestGuys/pifuxelck-server-go/server/db"
	"github.com/GreatestGuys/pifuxelck-server-go/server/handlers"
	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
	"github.com/gorilla/mux"
)

// Config defines all the options that can be configured for a running instance
// of the pifuxelck server.
type Config struct {
	Port        int
	DBConfig db.Config
}

// Run takes a Config and runs the pifuxelck server indefinitely.
func Run(config Config) {
	log.Infof("Starting pifuxelck server.")

	address := ":" + strconv.Itoa(config.Port)
	log.Infof("Listening on port %v.", config.Port)

	db.Init(config.DBConfig)

	http.Handle("/", newRouter())
	http.ListenAndServe(address, nil)
}

func newRouter() *mux.Router {
	r := mux.NewRouter()

	r.Handle("/secret", handlers.Secret)
	r.HandleFunc("/", handlers.Home)

	return r
}
