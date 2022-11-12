package server

import (
	k8sclient "moodle/client/k8sclient"
	mariaclient "moodle/client/mariaclient"
	"moodle/database"
)
	


type dep struct {
	MongoDB database.DBClient
	MariaClient mariaclient.MariaClient
	K8sClient k8sclient.K8sClient
	// Handler *Handler
}


func initDependencies() *dep {
	return &dep{
		MongoDB:     database.DBClient{},
		MariaClient: mariaclient.MariaClient{},
		K8sClient:   k8sclient.K8sClient{},
	}
}

