package worker

import (
	maria "moodle/internal/repository/moodle"
	mongo "moodle/pkg/mongodbiface"
	"moodle/config"
	nfs "moodle/pkg/client"
	nfs4 "moodle/pkg/client/nfs4"
)


type NFSWorker struct {
	Maria *maria.MariaDB
	Mongo mongo.DB
	NFSClient	 *nfs4.NfsClient
}

// new mariadb worker
func NewNFSWorker(mongo mongo.DB) *NFSWorker {
	mariaclient := maria.NewMariaDB(config.MARIAHOSTW, 3306, config.MARIAUSER, config.MARIAPASSWORD)
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
