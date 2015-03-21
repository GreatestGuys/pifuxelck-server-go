package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/GreatestGuys/pifuxelck-server-go/server"
	pifuxelckLog "github.com/GreatestGuys/pifuxelck-server-go/server/log"
)

var port = flag.Int("port", 3000, "The port number to listen on.")

var logLevel = flag.Int("verbosity", 3,
	"The verbosity of the log statements. The larger the number, the more verbose.")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	log.SetFlags(log.Ldate | log.Ltime)
	pifuxelckLog.SetLogLevel(*logLevel)

	server.Run(server.Config{
		Port: *port,
	})
}
