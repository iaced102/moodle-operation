package client

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
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


type K8sClient struct {
}

// new k8s
func NewK8sClient() *K8sClient {
	return &K8sClient{}
}

func (client *K8sClient) NewClientSet() *kubernetes.Clientset {
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
func (client *K8sClient) ListNamespaces(clientset *kubernetes.Clientset) []string {
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
func (client *K8sClient) ListPods(clientset *kubernetes.Clientset, namespace string) []string {
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
func (client *K8sClient) ListStorageClasses(clientset *kubernetes.Clientset) []string {
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
func (client *K8sClient) SetDefaultStorageClass(clientset *kubernetes.Clientset, storageClass string) {
	fmt.Printf("Setting default storage class to %q: ", storageClass)
	_, err := clientset.StorageV1().StorageClasses().Patch(context.Background(), storageClass, types.JSONPatchType, []byte(`[{"op": "replace", "path": "/metadata/annotations/storageclass.kubernetes.io~1is-default-class", "value": "true"}]`), metav1.PatchOptions{})
	if err != nil {
		panic(err.Error())
	}
}


// mapping env vars to values
func (client *K8sClient) GetEnvVarsValues(clientset *kubernetes.Clientset, namespace string, pod string) map[string]string {
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
func (client *K8sClient) CreateNamespace(clientset *kubernetes.Clientset, namespace string) {
	fmt.Printf("Creating namespace %q\n", namespace)
	_, err := clientset.CoreV1().Namespaces().Create(context.Background(), &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
}


// delete namespace
func (client *K8sClient) DeleteNamespace(clientset *kubernetes.Clientset, namespace string) {
	fmt.Printf("Deleting namespace %q:\n", namespace)
	err := clientset.CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
}


// get service external IP
func (client *K8sClient) GetServiceExternalIP(clientset *kubernetes.Clientset, namespace string, service string) string {
	fmt.Printf("Getting service external IP from service %q in namespace %q:\n", service, namespace)
	svc, err := clientset.CoreV1().Services(namespace).Get(context.Background(), service, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return svc.Status.LoadBalancer.Ingress[0].IP
}



// apply pvc from filepath
func (client *K8sClient) ApplyPVC(clientset *kubernetes.Clientset, namespace string, filepath string) {
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
func (client *K8sClient) ApplyService(clientset *kubernetes.Clientset, namespace string, filepath string) {
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
func (client *K8sClient) ApplyStatefulSet(clientset *kubernetes.Clientset, namespace string, filepath string) {
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
func (client *K8sClient) DeletePVC(clientset *kubernetes.Clientset, namespace string, pvc string) {
	fmt.Printf("Deleting pvc %q in namespace %q:\n", pvc, namespace)
	err := clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(context.Background(), pvc, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
}


// delete service
func (client *K8sClient) DeleteService(clientset *kubernetes.Clientset, namespace string, service string) {
	fmt.Printf("Deleting service %q in namespace %q:\n", service, namespace)
	err := clientset.CoreV1().Services(namespace).Delete(context.Background(), service, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
}


// delete statefulset
func (client *K8sClient) DeleteStatefulSet(clientset *kubernetes.Clientset, namespace string, statefulset string) {
	fmt.Printf("Deleting statefulset %q in namespace %q:\n", statefulset, namespace)
	err := clientset.AppsV1().StatefulSets(namespace).Delete(context.Background(), statefulset, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
}


// get pod logs
func (client *K8sClient) GetPodLogs(clientset *kubernetes.Clientset, namespace string, pod string) {
	fmt.Printf("Getting pod logs %q in namespace %q:\n", pod, namespace)
	podLogOpts := v1.PodLogOptions{}
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
	fmt.Printf("%s", str)
}


// return statefulset list
func (client *K8sClient) ListStatefulSet(clientset *kubernetes.Clientset, namespace string) []string {
	fmt.Printf("Listing statefulset in namespace %q:\n", namespace)
	statefulsetList, err := clientset.AppsV1().StatefulSets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	var statefulset []string
	for _, s := range statefulsetList.Items {
		statefulset = append(statefulset, s.Name)
	}
	return statefulset
}


// return services list
func (client *K8sClient) ListService(clientset *kubernetes.Clientset, namespace string) []string {
	fmt.Printf("Listing services in namespace %q:\n", namespace)
	serviceList, err := clientset.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	var service []string
	for _, s := range serviceList.Items {
		service = append(service, s.Name)
	}
	return service
}


// return pvc list
func (client *K8sClient) ListPVC(clientset *kubernetes.Clientset, namespace string) []string {
	fmt.Printf("Listing pvc in namespace %q:\n", namespace)
	pvcs, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	var pvcSlice []string
	for _, pvc := range pvcs.Items {
		fmt.Printf("PVC: %q", pvc.Name)
		pvcSlice = append(pvcSlice, pvc.Name)
	}
	return pvcSlice
}
