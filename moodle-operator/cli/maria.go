package cli

import (
	"fmt"
	client "moodle/client/mariaclient"
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
	default:
		fmt.Println("Invalid command")
	}
}
