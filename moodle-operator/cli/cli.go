package cli

import (
	"fmt"
	"moodle/k8sclient"
	"os"
)

type CLI struct{}

func (cli *CLI) Run() {
	client := k8sclient.K8S{}
	switch os.Args[1] {
	case "list-namespaces":
		fmt.Println(client.ListNamespaces(client.NewClientSet()))
	case "list-pods":
		fmt.Println(client.ListPods(client.NewClientSet(), os.Args[2]))
	case "list-storage-classes":
		fmt.Println(client.ListStorageClasses(client.NewClientSet()))
	case "set-default-storage-class":
		client.SetDefaultStorageClass(client.NewClientSet(), os.Args[2])
	case "get-env-vars-values":
		fmt.Println(client.GetEnvVarsValues(client.NewClientSet(), os.Args[2], os.Args[3]))
	case "create-namespace":
		client.CreateNamespace(client.NewClientSet(), os.Args[2])
	case "delete-namespace":
		client.DeleteNamespace(client.NewClientSet(), os.Args[2])
	case "get-service":
		fmt.Println(client.GetServiceExternalIP(client.NewClientSet(), os.Args[2], os.Args[3]))
	case "apply-pvc":
		client.ApplyPVC(client.NewClientSet(), os.Args[2], os.Args[3])
	case "apply-service":
		client.ApplyService(client.NewClientSet(), os.Args[2], os.Args[3])
	case "apply-statefulset":
		client.ApplyStatefulSet(client.NewClientSet(), os.Args[2], os.Args[3])
	case "delete-pvc":
		client.DeletePVC(client.NewClientSet(), os.Args[2], os.Args[3])
	case "delete-service":
		client.DeleteService(client.NewClientSet(), os.Args[2], os.Args[3])
	case "delete-statefulset":
		client.DeleteStatefulSet(client.NewClientSet(), os.Args[2], os.Args[3])
	case "list-pvc":
		fmt.Println(client.ListPVC(client.NewClientSet(), os.Args[2]))

	default:
		fmt.Println("Invalid command")
	}
}
