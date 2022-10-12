#!/bin/bsh
helm='/usr/bin/helm'


install_nfs_provisioner() {
	nfs_name='nfs-subdir-external-provisioner'
	nfs_package='nfs-subdir-external-provisioner/nfs-subdir-external-provisioner'
	nfs_ip='123.30.234.224'
	nfs_path='/moodle'
	$helm repo add nfs-subdir-external-provisioner https://kubernetes-sigs.github.io/nfs-subdir-external-provisioner/ 
	$helm install $nfs_name $nfs_package --set nfs.server=$nfs_ip  --set nfs.path=$nfs_path
}


get_storage_class() {
	kubectl get storageclass
}


toggle_default_storage_class() {
	premium_ssd1=${1:-'false'}
	nfs=${2:-'true'}
	# kubectl patch storageclass standard -p "{\"metadata\": {\"annotations\":{\"storageclass.kubernetes.io/is-default-class\":\"$standard\"}}}"
	kubectl patch storageclass premium-ssd1  -p "{\"metadata\": {\"annotations\":{\"storageclass.kubernetes.io/is-default-class\":\"$premium_ssd1\"}}}"
	kubectl patch storageclass nfs-client -p "{\"metadata\": {\"annotations\":{\"storageclass.kubernetes.io/is-default-class\":\"$nfs\"}}}"
	echo -e "premium-ssd1: $premium_ssd1\nnfs: $nfs"
	# kubectl patch storageclass premium-ssd1 -p "{\"metadata\": {\"annotations\":{\"storageclass.kubernetes.io/is-default-class\":\"true\"}}}"
}


get_moodle_password() {
	ns=${1:-'my-release'}
	pod=`get_pod_name $ns`
	# name="${ns}"
	# echo Password: $(kubectl get secret --namespace  $ns $name -o jsonpath="{.data.moodle-password}" | base64 -d)
	kubectl -n $ns exec $pod -- printenv  | grep -i moodle_password 
}

get_mariadb_password() {
	ns=${1:-'my-release'}
	name="${ns}-mariadb"
	echo Password: $(kubectl get secret --namespace  $ns $name -o jsonpath="{.data.mariadb-root-password}" | base64 -d)
}

create_moodle() {
	# kubectl patch storageclass premium-ssd1 -p "{\"metadata\": {\"annotations\":{\"storageclass.kubernetes.io/is-default-class\":\"true\"}}}"
	name=${1:-'my-release'}
	ns=${name:-'my-release'}
	kubectl create ns $ns
	package=${2:-'bitnami/moodle'}
	$helm repo add bitnami https://charts.bitnami.com/bitnami
	$helm  install $name $package --namespace $ns -f values.yaml
}


delete_moodle() {
	name=${1:-'my-release'}
	ns=${name:-'my-release'}
	$helm delete $name  --namespace $ns
	kubectl delete pvc "data-${name}-mariadb-0" -n $ns
	pv=`kubectl get pv -n $ns | awk '/pvc/ {print $1}'`
	kubectl delete pv $pv -n $ns
	kubectl delete ns $ns
}

forward_port() {
	host=${1:-'localhost'}
	service_name=${2:-'my-release-moodle'}
	kubectl port-forward --address $host services/$service_name 8080:80
}


list_url() {
	minikube service list
}



get_pods() {
	ns=${1:-'my-release'}
	kubectl get pods -n $ns
}


get_pod_name() {
	ns=${1:-'my-release'}
	kubectl get pods -n $ns | awk '/moodle/ {print $1}' | head -n 1
}

get_maria_pod_name() {
	ns=${1:-'my-release'}
	kubectl get pods -n $ns | awk '/maria/ {print $1}' | head -n 1
}

get_pvc() {
	ns=${1:-'my-release'}
	kubectl get pvc -n $ns
}


get_pv() {
	ns=${1:-'my-release'}
	kubectl get pvc -n $ns
}


get_events() {
	ns=${1:-'my-release'}
	kubectl get events --sort-by='.lastTimestamp' -n $ns
}


get_ns() {
	kubectl get ns
}


get_svc() {
	ns=${1:-'my-release'}
	kubectl get svc -n $ns
}


get_php_version() {
	ns=${1:-'my-release'}
	pod_name=`get_pod_name $ns`
	echo `kubectl  -n $ns exec -it $pod_name -- /opt/bitnami/php/bin/php -v | awk -F' ' '/^PHP/ {print $2}'`
}

copy_to_pod() {
	ns=${1:-'my-release'}
	pod_name=$2
	kubectl cp $3  ${ns}/${pod_name}:/tmp
}

copy_from_pod() {
	ns=${1:-'my-release'}
	pod_name=$2
	kubectl cp ${ns}/${pod_name}:/tmp/usecase1.sql $3
}

kubectl_exec() {
	ns=${1:-'my-release'}
	pod_name=$2
	kubectl -n $ns exec -i $pod_name -- bash -c "$3 $4 $5 $6 $7 $8 $9 $10 $11 $12"
}


install_plugin() {
	echo 'hele'
}


restore_moodle_db() {
	ns=${1:-'my-release'}
	pod=`get_maria_pod_name $ns`
	copy_to_pod  $ns $pod ./usecase1.sql
	# mysqldump -u root -plCjADA6RV8 bitnami_moodle > /tmp/usecase1.sql
	kubectl -n newtheme exec -i newtheme-mariadb-0 -- bash -c "/opt/bitnami/mariadb/bin/mysql -u root -plCjADA6RV8 -D bitnami_moodle < /tmp/usecase1.sql"

}

main() {
	read -p 'enter your cluster name: ' clustername
	read -p 'enter your site name: ' sitename
	read -p 'enter your ccu: ' ccu

	moodle_version_available=('4.0.4-debian-11-r4' '4.0.3-debian-11-r7' \
		'4.0.2-debian-11-r13' '4.0.1-debian-11-r12' '3.11.10-debian-11-r4' \
		'3.10.4-debian-10-r6' '3.9.2-debian-10-r47' '3.8.3' '3.7.3-ol-7-r8' \
		'3.6.4-r6' '3.5.3' '3.4.2' '3.3.2-r1' \
	)
	echo -e "moodle version available is below:\n"
	for i in {0..12}; do echo "$i: ${moodle_version_available[$i]}"; done;
	read -p "choose moodle version you want: " _index
	# moodle_version=${moodle_version_available[$_index]}
	sed -i "86s/moodleSiteName.*$/moodleSiteName: \"$sitename\"/g" ./values.yaml
	moodle_version='4.0.4-debian-11-r4'
	sed -i "58s/tag.*$/tag: \"$moodle_version\"/g" ./values.yaml

	create_moodle $clustername

}

$*
