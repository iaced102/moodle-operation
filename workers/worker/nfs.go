package worker

import (
	server "moodle/cmd/restserver/server"
	maria "moodle/internal/repository/moodle"
	nfs "moodle/pkg/client"
	nfs4 "moodle/pkg/client/nfs4"
	mongo "moodle/pkg/mongodbiface"
)


type NFSWorker struct {
	Maria *maria.MariaDB
	Mongo mongo.DB
	NFSClient	 *nfs4.NfsClient
}

// new mariadb worker
func NewNFSWorker(mongo mongo.DB) *NFSWorker {
	mariaclient := maria.NewMariaDB(server.NewMariaDB())
	nfsclient := nfs.NewNFSClient()
	
	return &NFSWorker{
		Maria: mariaclient,
		Mongo: mongo,
		NFSClient: nfsclient,
	}
}

// update lms logo by upload file to nfs server
func (w *NFSWorker) UpdateLmsLogo() {
}


// update lms favicon by upload file to nfs server
