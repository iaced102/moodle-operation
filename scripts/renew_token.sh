#!/bin/bash


password="MzU2ZmE3ZDM5MmUwMDFkYThiMzkwMmE3"
username="duynn@bizflycloud.vn"

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
	curl -s 'https://manage.bizflycloud.vn/api/kubernetes-engine/_/6o0cn9lv42livqek/kubeconfig' \
	  -H 'authority: manage.bizflycloud.vn' \
	  -H 'accept: application/json, text/plain, */*' \
	  -H "X-Tenant-Name: duynn@bizflycloud.vn" \
	  -H "X-Auth-Token: $token" 
}

token=$(get_token)
mongosh "mongodb://localhost:27017/moodle" --eval "db.tokens.updateMany({\"id\": \"1\"}, {\$set: {\"token\": \"$token\"}})"

