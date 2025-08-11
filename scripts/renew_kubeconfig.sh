#!/bin/bash


password="47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU"
username="phucvuhuy@vccorp.vn"

function get_token() {
    curl -s --request POST 'https://manage.bizflycloud.vn/api/token' \
        --header 'Content-Type: application/json' \
        -d "{
            \"username\": \"$username\",
            \"password\": \"$password\",
            \"auth_method\": \"password\"
        }"  | jq -r .token
}


function get_kubeconfig() {
	token=$(get_token)
	curl -s 'https://manage.bizflycloud.vn/api/kubernetes-engine/_/hnheekfmpirevw4y/kubeconfig?expire_time=525948' \
	  -H 'authority: manage.bizflycloud.vn' \
	  -H 'accept: application/json, text/plain, */*' \
	  -H "X-Tenant-Name: bizflycloud@vccloud.vn" \
	  -H "X-Auth-Token: $token" 
}

get_kubeconfig > $HOME/moodle-cluster.kubeconfig
