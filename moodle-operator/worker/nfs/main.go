package worker

import (
	mariaAdapter "moodle/adapter/mariadb"
	mongoAdapter "moodle/adapter/mongo"
	nfs "moodle/client/nfs"
	nfs4 "moodle/client/nfs/nfs4"
	"moodle/config"
)


type NFSWorker struct {
	MariaAdapter *mariaAdapter.MariaAdapter
	MongoAdapter mongoAdapter.MongoAdapter
	NFSClient	 *nfs4.NfsClient
}

// new mariadb worker
func NewMariaWorker(adapter mongoAdapter.MongoAdapter) *NFSWorker {
	mariaclient := mariaAdapter.NewMariaAdapter(config.MARIAHOSTW, 3306, config.MARIAUSER, config.MARIAPASSWORD)
	nfsclient := nfs.NewNFSClient()
	
	return &NFSWorker{
		MariaAdapter: mariaclient,
		MongoAdapter: adapter,
		NFSClient: nfsclient,
	}
}

// update lms logo by upload file to nfs server
func (w *NFSWorker) UpdateLmsLogo() {
}


// update lms favicon by upload file to nfs server
