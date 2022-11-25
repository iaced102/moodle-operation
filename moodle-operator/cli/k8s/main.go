package cli

import (
	"fmt"
	client "moodle/client/k8s"
	"os"
)

type K8sCli struct{}

func NewK8sCli() *K8sCli {
	return &K8sCli{}
}

func (cli *K8sCli) Run() {
	// k8s 
	k8s := client.NewK8sClient()
	clientset := k8s.NewClientSet()
	switch os.Args[1] {
	case "list-namespaces":
		fmt.Println(k8s.ListNamespaces(clientset))
	case "list-pods":
		fmt.Println(k8s.ListPods(clientset, os.Args[2]))
	case "list-storage-classes":
		fmt.Println(k8s.ListStorageClasses(clientset))
	case "set-default-storage-class":
		k8s.SetDefaultStorageClass(clientset, os.Args[2])
	case "get-env-vars-values":
		fmt.Println(k8s.GetEnvVarsValues(clientset, os.Args[2], os.Args[3]))
	case "create-namespace":
		k8s.CreateNamespace(clientset, os.Args[2])
	case "delete-namespace":
		k8s.DeleteNamespace(clientset, os.Args[2])
	case "get-service":
		fmt.Println(k8s.GetServiceExternalIP(clientset, os.Args[2], os.Args[3]))
	case "apply-pvc":
		k8s.ApplyPVC(clientset, os.Args[2])
	case "apply-service":
		k8s.ApplyService(clientset, os.Args[2])
	case "apply-statefulset":
		k8s.ApplyStatefulSet(clientset, os.Args[2], os.Args[3])
	case "delete-pvc":
		k8s.DeletePVC(clientset, os.Args[2], os.Args[3])
	case "delete-service":
		k8s.DeleteService(clientset, os.Args[2], os.Args[3])
	case "delete-statefulset":
		k8s.DeleteStatefulSet(clientset, os.Args[2], os.Args[3])
	case "list-pvc":
		fmt.Println(k8s.ListPVC(clientset, os.Args[2]))
	case "get-podlogs":
		k8s.GetPodLogs(clientset, os.Args[2], os.Args[3])
	case "list-statefulsets":
		fmt.Println(k8s.ListStatefulSet(clientset, os.Args[2]))
	case "list-services":
		fmt.Println(k8s.ListService(clientset, os.Args[2]))
	default:
		fmt.Println("Invalid command")
	}
}
