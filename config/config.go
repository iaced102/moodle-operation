package config

import "k8s.io/client-go/util/homedir"


var (
	HOME = homedir.HomeDir()
	MONGOURI = "mongodb://localhost:27017"
	DBNAME = "moodle"
	USERNAME = "duynn@bizflycloud.vn"
	INTERVAL = 10
	MOODLE_TOKEN = "gAAAAABjvhxZwCj1n8fFu-nFALaMr4KKZ7G3UPI7pAETzVxGa4pYuWdUelC6DcjLSRd8rHXIWVsya4q0qijXmOiwPbJzfNOSxBthvRZmGBxfc46rHefbWpwXq-EqFAsxO5MDFPYjUsBocTwCtb6UFdP8ruEdAWxWQ9YQl2KPAuyP_bprvmZwUkI"
	USER = "admin"
	MARIAHOSTR = "45.124.94.92"
	MARIAHOSTW = "45.124.94.92"
	MARIAUSER = "root"
	MARIAPASSWORD = "wLzTiBkkhCsenUXHpCYLqQ5pNvuLGTuMUTsi"
	MARIAPORT = 3306
	NFSSEVER = "103.148.57.178:2049"
	PVC_FILEPATH = HOME + "/gits/moodle-operator/deploy/moodle/pvc.yaml"
	SECRET_FILEPATH = HOME + "/gits/moodle-operator/deploy/moodle/secret.yaml"
	SERVICE_FILEPATH = HOME + "/gits/moodle-operator/deploy/moodle/service.yaml"
	STATEFULSET_FILEPATH = HOME + "/gits/moodle-operator/deploy/moodle/statefulset.yaml"
	INGRESS_FILEPATH = HOME + "/gits/moodle-operator/deploy/moodle/ingress.yaml"
	KUBECONFIG = "moodle-cluster.kubeconfig"
	CLUSTERID = "6o0cn9lv42livqek"
	SENDMAILDOMAIN = "https://sendmail.bizflycloud.vn"
	LMSDOMAIN = "lms.bfcplatform.vn"
)
