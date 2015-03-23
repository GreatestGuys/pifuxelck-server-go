package server

import (
	"net/http"
	"strconv"

	"github.com/GreatestGuys/pifuxelck-server-go/server/db"
	"github.com/GreatestGuys/pifuxelck-server-go/server/handlers"
	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

// Config defines all the options that can be configured for a running instance
// of the pifuxelck server.
type Config struct {
	Port     int
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

	// Install the prometheus handler at /metrics. This will allow the prometheus
	// monitoring system to build time series of the servers status.
	r.Handle("/metrics", prometheus.Handler())

	s := r.PathPrefix("/api/2/").Subrouter()
	handlers.InstallAccountHandlers(s)
	handlers.InstallContactHandlers(s)

	return r
}
