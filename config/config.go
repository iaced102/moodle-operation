package config
import "k8s.io/client-go/util/homedir"


var (
	HOME = homedir.HomeDir()
	MONGOURI = "mongodb://localhost:27017"
	DBNAME = "moodle"
	USERNAME = "duynn@bizflycloud.vn"
	PASSWORD = "MzU2ZmE3ZDM5MmUwMDFkYThiMzkwMmE3"
	BIZFLYCLOUD_TOKEN = "gAAAAABjxQHVqjvryJlPvrm15_OdlNDk6RxhYwl291ivFlP78rdgLCv4M7HgX8IB2hUZV3Ge-qBVvFNl37nnptxsLoIs_uWTLshlLO__woefRpfwydU23oUF4Oxj1iT7ZsxgyOP278D4As6g35qG-geGhgdkEGZhgsythyuWh44qdz9ux22bK5U"
	INTERVAL = 10
	MOODLE_TOKEN = "gAAAAABjvhxZwCj1n8fFu-nFALaMr4KKZ7G3UPI7pAETzVxGa4pYuWdUelC6DcjLSRd8rHXIWVsya4q0qijXmOiwPbJzfNOSxBthvRZmGBxfc46rHefbWpwXq-EqFAsxO5MDFPYjUsBocTwCtb6UFdP8ruEdAWxWQ9YQl2KPAuyP_bprvmZwUkI"
	USER = "admin"
	MARIAHOSTR = "45.124.94.39"
	MARIAHOSTW = "45.124.94.112"
	MARIAUSER = "root"
	MARIAPASSWORD = "0YU8381WUlk1u9ysVbF4Qb5FigNW8z8uCvPI"
	MARIAPORT = 3306
	NFSSEVER = "123.30.234.224:2049"
	PVC_FILEPATH = HOME + "/gits/moodle-operator/deploy/moodle/pvc.yaml"
	SERVICE_FILEPATH = HOME + "/gits/moodle-operator/deploy/moodle/service.yaml"
	STATEFULSET_FILEPATH = HOME + "/gits/moodle-operator/deploy/moodle/statefulset.yaml"
	INGRESS_FILEPATH = HOME + "/gits/moodle-operator/deploy/moodle/ingress.yaml"
	KUBECONFIG = "moodle-cluster.kubeconfig"
	CLUSTERID = "6o0cn9lv42livqek"
	SENDMAILDOMAIN = "https://sendmail.bizflycloud.vn"
	LMSDOMAIN = "lms.bfcplatform.vn"
)
