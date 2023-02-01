package client

import (
	"context"
	"fmt"
	"moodle/client/nfs/nfs4"
	"moodle/config"
)

func NewNFSClient() *nfs4.NfsClient {
	nfsclient , err := nfs4.NewNfsClient(context.Background(), config.NFSSEVER,
		nfs4.AuthParams{
			Uid: 0,
			Gid: 0,
			MachineName: "nfs",
		})
	if err != nil {
		fmt.Println(err)
		return  nil
	}
	return nfsclient
}
