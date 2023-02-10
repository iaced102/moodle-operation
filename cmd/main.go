package main

import (
	"moodle/cmd/restserver/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	go server.RunWorker()
	server.Start()
}

