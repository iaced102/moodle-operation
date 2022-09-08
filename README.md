![moodle_operator_CI_Prod](https://github.com/Netzlink/moodle-operator/workflows/moodle_operator_CI_Prod/badge.svg?branch=master)
# moodle Operator ( working on OpenShift)
This is an operator for moodle on K8s build with OpenShift in mind.

# Deployment
If you use OpenShift just change kubectl to oc.
```bash
git clone git@github.com:Netzlink/moodle-operator.git
cd moodle-operator
kubectl apply -f moodle-operator/deploy/service_account.yaml
kubectl apply -f moodle-operator/deploy/role.yaml
kubectl apply -f moodle-operator/deploy/role_binding.yaml 
kubectl apply -f moodle-operator/deploy/crds/moodle.netzlink.com_moodles_crd.yaml
kubectl apply -f moodle-operator/deploy/operator.yaml
kubectl apply -f moodle-operator/deploy/crds/moodle.netzlink.com_v1alpha1_moodle_cr.yaml
```
verify it's woriking
```bash
kubectl get pods -w
```
you should be ready to go!
# Todo
- For now the default systemaccount in the namespace has to have the anyuid SCC in OpenShift
# Thanks
- RedHat for the awesome operator-framework
- ASO for introducing me to it
- Bitnami for there Helm file for moodle
- CNCF for K8S 
