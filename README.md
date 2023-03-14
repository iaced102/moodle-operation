export PVC_FILEPATH=$HOME/gits/moodle-operator/deploy/moodle-default/pvc.yaml
export SERVICE_FILEPATH=$HOME/gits/moodle-operator/deploy/moodle-default/service.yaml
export STATEFULSET_FILEPATH=$HOME/gits/moodle-operator/deploy/moodle
export DEPLOYMENT_FILEPATH=$HOME/gits/moodle-operator/deploy/moodle
export KUBECONFIG=$HOME/moodle-cluster.kubeconfig

# API
```
go run main.go
```

# CLI
```
go run main.go COMMAND
```

# filedir 
## logo: 
### path : /2b/83/2b831e652219e350659e4a71af9c4b9c4c99411c
### url: http://hello1.lms.bizflycloud.vn/pluginfile.php/1/theme_edumy/headerlogo1/1676573021/Logo mới Bizfly Cloud-01.png

## favicon:  72x72 jpeg required
### path : /35/1d/351d302a3d094fdde4638e7dbdde86af3199be7e
### url: http://hello.lms.bizflycloud.vn/pluginfile.php/1/theme_edumy/favicon/1676573021/z3665638475480_d6dab64f97b26f1c73cd539411ae9990.jpg



# adapt db before restore from backup files
```
mysqldump -u root -p moodle > moodle01112022.sql
mysql -u duy -p5Yk7741J2JVWPTQkT9eKkcbAaTUs5XzTvIFL moodle < moodle01112022
sed -i 's/utf8mb4_0900_ai_ci/utf8_unicode_ci/g' moodle.sql
sed -i 's/utf8mb4/utf8/g' moodle.sql
sed -i 's/utf8_unicode_520_ci/utf8_unicode_ci/g' moodle.sql
```



# install certmanger to use let's encrypt
## install
```
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.5.4/cert-manager.yaml
```
## verify
```
kubectl get pods --namespace cert-manager
```
## Create a ClusterIssuer resource to configure cert-manager with your Let's Encrypt account credentials.
```clusterissuer.yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    email: your-email-address
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
```
