package k8sclient

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	appsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/yaml"
)


type K8S struct {
}


func (client *K8S) NewClientSet() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, "moodle.kubeconfig"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("/home/duy/moodle.kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

// return all namespaces
func (client *K8S) ListNamespaces(clientset *kubernetes.Clientset) []string {
	fmt.Println("Listing namespaces:")
	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	var ns []string
	for _, namespace := range namespaces.Items {
		ns = append(ns, namespace.ObjectMeta.Name)
	}
	return ns
}


// list pods
func (client *K8S) ListPods(clientset *kubernetes.Clientset, namespace string) []string {
	fmt.Printf("Listing pods in namespace %q:\n", namespace)
	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	var podsList []string
	for _, pod := range pods.Items {
		podsList = append(podsList, pod.ObjectMeta.Name)
	}
	return podsList
}


// list storage classes
func (client *K8S) ListStorageClasses(clientset *kubernetes.Clientset) []string {
	fmt.Println("Listing storage classes:")
	storageClasses, err := clientset.StorageV1().StorageClasses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	var sc []string
	for _, storageClass := range storageClasses.Items {
		sc = append(sc, storageClass.ObjectMeta.Name)
	}
	return sc
}


// set default storage class
func (client *K8S) SetDefaultStorageClass(clientset *kubernetes.Clientset, storageClass string) {
	fmt.Printf("Setting default storage class to %q: ", storageClass)
	_, err := clientset.StorageV1().StorageClasses().Patch(context.Background(), storageClass, types.JSONPatchType, []byte(`[{"op": "replace", "path": "/metadata/annotations/storageclass.kubernetes.io~1is-default-class", "value": "true"}]`), metav1.PatchOptions{})
	if err != nil {
		panic(err.Error())
	}
}


// mapping env vars to values
func (client *K8S) GetEnvVarsValues(clientset *kubernetes.Clientset, namespace string, pod string) map[string]string {
	fmt.Printf("Getting env vars from pod %q in namespace %q:\n", pod, namespace)
	env, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), pod, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	envVars := make(map[string]string)
	for _, envVar := range env.Spec.Containers[0].Env {
		envVars[envVar.Name] = envVar.Value
	}
	return envVars
}


// create namespace
func (client *K8S) CreateNamespace(clientset *kubernetes.Clientset, namespace string) {
	fmt.Printf("Creating namespace %q\n", namespace)
	_, err := clientset.CoreV1().Namespaces().Create(context.Background(), &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
}


// delete namespace
func (client *K8S) DeleteNamespace(clientset *kubernetes.Clientset, namespace string) {
	fmt.Printf("Deleting namespace %q:\n", namespace)
	err := clientset.CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
}


// get service external IP
func (client *K8S) GetServiceExternalIP(clientset *kubernetes.Clientset, namespace string, service string) string {
	fmt.Printf("Getting service external IP from service %q in namespace %q:\n", service, namespace)
	svc, err := clientset.CoreV1().Services(namespace).Get(context.Background(), service, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return svc.Status.LoadBalancer.Ingress[0].IP
}



// apply pvc from filepath
func (client *K8S) ApplyPVC(clientset *kubernetes.Clientset, namespace string, filepath string) {
	fmt.Printf("Deploying pvc from file %q in namespace %q:\n", filepath, namespace)
	// read file
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err.Error())
	}
	// unmarshal file
	var pvc corev1.PersistentVolumeClaimApplyConfiguration
	err = yaml.Unmarshal(file, &pvc)
	if err != nil {
		panic(err.Error())
	}
	// apply pvc
	_, err = clientset.CoreV1().PersistentVolumeClaims(namespace).Apply(context.Background(), &pvc, metav1.ApplyOptions{FieldManager: "kubectl-client-side-apply"})
	if err != nil {
		panic(err.Error())
	}
}

// apply service from filepath
func (client *K8S) ApplyService(clientset *kubernetes.Clientset, namespace string, filepath string) {
	fmt.Printf("Applying service from file %q in namespace %q:\n", filepath, namespace)
	// read file
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err.Error())
	}
	// unmarshal file
	var service corev1.ServiceApplyConfiguration
	err = yaml.Unmarshal(file, &service)
	if err != nil {
		panic(err.Error())
	}
	// apply service
	_, err = clientset.CoreV1().Services(namespace).Apply(context.Background(), &service, metav1.ApplyOptions{FieldManager: "kubectl-client-side-apply"})
	if err != nil {
		panic(err.Error())
	}
}


//  apply statefulset from  filepath
func (client *K8S) ApplyStatefulSet(clientset *kubernetes.Clientset, namespace string, filepath string) {
	fmt.Printf("Applying statefulset from file %q in namespace %q:\n", filepath, namespace)
	// read file
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err.Error())
	}
	// unmarshal file
	var statefulset appsv1.StatefulSetApplyConfiguration
	err = yaml.Unmarshal(file, &statefulset)
	if err != nil {
		panic(err.Error())
	}
	// apply statefulset
	_, err = clientset.AppsV1().StatefulSets(namespace).Apply(context.Background(), &statefulset, metav1.ApplyOptions{FieldManager: "kubectl-client-side-apply"})
	if err != nil {
		panic(err.Error())
	}
}


// delete pvc
func (client *K8S) DeletePVC(clientset *kubernetes.Clientset, namespace string, pvc string) {
	fmt.Printf("Deleting pvc %q in namespace %q:\n", pvc, namespace)
	err := clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(context.Background(), pvc, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
}


// delete service
func (client *K8S) DeleteService(clientset *kubernetes.Clientset, namespace string, service string) {
	fmt.Printf("Deleting service %q in namespace %q:\n", service, namespace)
	err := clientset.CoreV1().Services(namespace).Delete(context.Background(), service, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
}


// delete statefulset
func (client *K8S) DeleteStatefulSet(clientset *kubernetes.Clientset, namespace string, statefulset string) {
	fmt.Printf("Deleting statefulset %q in namespace %q:\n", statefulset, namespace)
	err := clientset.AppsV1().StatefulSets(namespace).Delete(context.Background(), statefulset, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
}


// get pod logs
func (client *K8S) GetPodLogs(clientset *kubernetes.Clientset, namespace string, pod string) {
	fmt.Printf("Getting pod logs %q in namespace %q:\n", pod, namespace)
	podLogOpts := corev1.PodLogOptions{}
	req := clientset.CoreV1().Pods(namespace).GetLogs(pod, &podLogOpts)
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		panic(err.Error())
	}
	defer podLogs.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		panic(err.Error())
	}
	str := buf.String()
	fmt.Printf("%s

", str)
}


// get events
func (client *K8S) GetEvents(clientset *kubernetes.Clientset, namespace string) {
	fmt.Printf("Getting events in namespace %q:\n", namespace)
	events, err := clientset.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, event := range events.Items {
		fmt.Printf("Event: %q

", event.Name)
	}
}


// lis statefulset
func (client *K8S) ListStatefulSet(clientset *kubernetes.Clientset, namespace string) {
	fmt.Printf("Listing statefulset in namespace %q:\n", namespace)
	statefulsets, err := clientset.AppsV1().StatefulSets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, statefulset := range statefulsets.Items {
		fmt.Printf("Statefulset: %q

", statefulset.Name)
	}
}


// list service
func (client *K8S) ListService(clientset *kubernetes.Clientset, namespace string) {
	fmt.Printf("Listing service in namespace %q:\n", namespace)
	services, err := clientset.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, service := range services.Items {
		fmt.Printf("Service: %q

", service.Name)
	}
}


// list pvc to slice
func (client *K8S) ListPVC(clientset *kubernetes.Clientset, namespace string) []string {
	fmt.Printf("Listing pvc in namespace %q:\n", namespace)
	pvcs, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	var pvcSlice []string
	for _, pvc := range pvcs.Items {
		fmt.Printf("PVC: %q

", pvc.Name)
		pvcSlice = append(pvcSlice, pvc.Name)
	}
	return pvcSlice
}
