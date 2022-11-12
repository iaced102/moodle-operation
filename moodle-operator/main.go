package main

import (
	// "moodle/cli"

	"moodle/restserver/server"

	log "github.com/sirupsen/logrus"
)


func main() {
	// k8scli := cli.NewK8sCli()
	// k8scli.Run()
	log.SetFormatter(&log.JSONFormatter{})
	server.Start()
}
