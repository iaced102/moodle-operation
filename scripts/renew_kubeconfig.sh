#!/bin/bash


password="ck)/i10RU%M8"
username="bizflycloud@vccloud.vn"

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
	curl -s 'https://manage.bizflycloud.vn/api/kubernetes-engine/_/hnheekfmpirevw4y/kubeconfig' \
	  -H 'authority: manage.bizflycloud.vn' \
	  -H 'accept: application/json, text/plain, */*' \
	  -H "X-Tenant-Name: bizflycloud@vccloud.vn" \
	  -H "X-Auth-Token: $token" 
}

get_kubeconfig > $HOME/moodle-cluster.kubeconfig
