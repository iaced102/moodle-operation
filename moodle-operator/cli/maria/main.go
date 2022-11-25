package cli

import (
	"fmt"
	client "moodle/client/maria"
	"os"
)


type MariaCli struct{}

func NewMariaCli() *MariaCli {
	return &MariaCli{}
}

func (cli *MariaCli) Run() {
	mariaclient := client.NewMariaClient()
	switch os.Args[1] {
		case "list-instance":
			mariaclient.ListInstances()
		case "create-instance":
			mariaclient.CreateInstance(os.Args[2])
		case "delete-instance":
			mariaclient.DeleteInstance(os.Args[2])
		case "create-instance-from-backup":
			mariaclient.CreateInstanceFromBackup(os.Args[2], os.Args[3])
	default:
		fmt.Println("Invalid command")
	}
}
