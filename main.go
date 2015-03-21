package main

import (
	"flag"
	"runtime"

	"github.com/GreatestGuys/pifuxelck-server-go/server"
	"github.com/GreatestGuys/pifuxelck-server-go/server/db"
	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
)

var port = flag.Int("port", 3000, "The port number to listen on.")

var logLevel = flag.Int("verbosity", 3,
	"The verbosity of the log statements. The larger the number, the more verbose.")

var mysqlHost = flag.String("mysql-host", "localhost",
	"The host running the pifuxelck MySQL server.")

var mysqlPort = flag.Int("mysql-port", 3306,
	"The port the pifuxelck MySQL server is listening on.")

var mysqlDB = flag.String("mysql-db", "pifuxelck",
	"The MySQL database that contains the pifuxelck game data.")

var mysqlUser = flag.String("mysql-user", "",
	"The username to use when connecting to the pifuxelck MySQL server.")

var mysqlPassword = flag.String("mysql-password", "",
	"The password to use when connecting to the pifuxelck MySQL server.")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	log.Init()
	log.SetLogLevel(*logLevel)

	server.Run(server.Config{
		Port: *port,
		DBConfig: db.Config{
			Host:     *mysqlHost,
			Port:     *mysqlPort,
			DB:       *mysqlDB,
			User:     *mysqlUser,
			Password: *mysqlPassword,
		},
	})
}
