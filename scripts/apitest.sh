#!/bin/bash
token="gAAAAABjvhxZwCj1n8fFu-nFALaMr4KKZ7G3UPI7pAETzVxGa4pYuWdUelC6DcjLSRd8rHXIWVsya4q0qijXmOiwPbJzfNOSxBthvRZmGBxfc46rHefbWpwXq-EqFAsxO5MDFPYjUsBocTwCtb6UFdP8ruEdAWxWQ9YQl2KPAuyP_bprvmZwUkI"

# create moodle
function moodle-create() {
	curl -i -X POST "http://localhost:5000/api/v1/moodles" \
		-H "X-Tenant-Name: admin" \
		-H "X-Auth-Token: $token" \
		-H "Content-Type: application/json" \
		-d "{
			\"email\": \"$1\",
			\"website_name\": \"$2\",
			\"pre_installed_course\": $3,
			\"packages_name\": \"$4\",
			\"autoscale\": $5,
			\"documents_storage_extra\": $6
		}"
}

# list moodle
function moodle-list() {
	curl -H "X-Tenant-Name: admin" \
		 -H "X-Auth-Token: $token" \
		 -iX GET "http://lms-service:5000/api/v1/moodles?page=1&limit=10&email=$1"
}

# get moodle
function moodle-get() {
	curl -H "X-Tenant-Name: admin" \
		 -H "X-Auth-Token: $token" \
		 -iX GET "http://localhost:5000/api/v1/moodles?id=$1" 
}

# search moodle
function search() {
	curl -iX GET "http://localhost:5000/api/v1/moodles?search=$1&page=1&limit=10&email=vctest-qatest-chaupth5-staging@vccloud.vn" \
		 -H "X-Auth-Token: $token" \
		 -H "X-Tenant-Name: admin"
}

# change packages moodle
function moodle-change-packages() {
	curl -i -X PUT "http://localhost:5000/api/v1/packages" \
		 -H "X-Auth-Token: $token" \
		 -H "X-Tenant-Name: admin" \
		 -H "Content-Type: application/json" \
		 -d "{
			\"moodle_id\": \"$1\",
			\"packages_name\": \"$2\"
		}"
}

# change autoscale moodle
function moodle-change-autoscale() {
	curl -i -X PUT "http://localhost:5000/api/v1/autoscale" \
		 -H "X-Auth-Token: $token" \
		 -H "X-Tenant-Name: admin" \
		 -H "Content-Type: application/json" \
		 -d "{
		 	\"moodle_id\": \"$1\",
		 	\"autoscale\": $2
		 }"
}

# change documents_storage_extra moodle
function moodle-change-document-storage-extra() {
	curl -i -X PUT "http://localhost:5000/api/v1/document-storage-extra" \
		 -H "Content-Type: application/json" \
		 -H "X-Auth-Token: $token" \
		 -H "X-Tenant-Name: admin" \
		 -d "{
		 	\"moodle_id\": \"$1\",
		 	\"documents_storage_extra\": $2
		 }"
}

# change pre_installed_course moodle
function moodle-change-pre-installed-course() {
	curl -i -X PUT "http://localhost:5000/api/v1/pre-installed-course" \
		 -H "X-Auth-Token: $token" \
		 -H "X-Tenant-Name: admin" \
		 -H "Content-Type: application/json" \
		 -d "{
		 	\"moodle_id\": \"$1\",
		 	\"pre_installed_course\": $2
		 }"
}

# delete moodle
function moodle-delete() {
	curl -i -X DELETE "http://localhost:5000/api/v1/moodles" \
		 -H "Content-Type: application/json" \
		 -H "X-Auth-Token: $token" \
		 -H "X-Tenant-Name: admin" \
		 -d "{
		 	\"moodle_id\": \"$1\"
		 }"
}

# list course
function course-list() {
	curl -H "X-Tenant-Name: admin" \
		 -H "X-Auth-Token: $token" \
		 -iX GET "http://localhost:5000/api/v1/courses?page=&limit="
}
# list pacakge
function package-list() {
	curl -H "X-Tenant-Name: admin" \
		 -H "X-Auth-Token: $token" \
		 -iX GET "http://localhost:5000/api/v1/packages?page=1&limit=10"
}

# upload logo
function upload-logo() {
	curl -F "file=@/home/duy/Pictures/duydeptrai.jpg" "http://localhost:5000/api/v1/logo?moodle_id=d7ba7a3c-cb4e-4cf7-ab8a-f4eec8a16e33" \
		 -H "X-Tenant-Name: admin" \
		 -H "X-Auth-Token: $token" 
}
function upload-favicon() {
	curl -F "file=@/home/duy/Pictures/duydeptrai.jpg" -i "http://localhost:5000/api/v1/favicon?moodle_id=d7ba7a3c-cb4e-4cf7-ab8a-f4eec8a16e33" \
		 -H "X-Tenant-Name: admin" \
		 -H "X-Auth-Token: $token" 
}

function upload-video() {
	curl -X POST "http://localhost:5000/api/v1/video" \
		 -H "X-Tenant-Name: admin" \
		 -H "X-Auth-Token: $token" \
		 -d '{
		 	"moodle_id": "d7ba7a3c-cb4e-4cf7-ab8a-f4eec8a16e33",
		 	"video_url": "https://www.youtube.com/watch?v=QH2-TGUlwu4"
		 }'
}

function upload-vision-image() {
	curl -F "file=@/home/duy/Pictures/duydeptrai.jpg" "http://localhost:5000/api/v1/vision-image?moodle_id=d7ba7a3c-cb4e-4cf7-ab8a-f4eec8a16e33" \
		 -H "X-Tenant-Name: admin" \
		 -H "X-Auth-Token: $token" 
}

function upload-vision-content() {
	curl -X POST "http://localhost:5000/api/v1/vision-content" \
		 -H "X-Tenant-Name: admin" \
		 -H "X-Auth-Token: $token" \
		 -d '{
		 	"moodle_id": "d7ba7a3c-cb4e-4cf7-ab8a-f4eec8a16e33",
		 	"title": "test",
		 	"content": "https://www.youtube.com/watch?v=QH2-TGUlwu4"
		 }'
}

function upload-banner() {
	curl -F "file=@/home/duy/Pictures/duydeptrai.jpg"  "http://localhost:5000/api/v1/banner?moodle_id=d7ba7a3c-cb4e-4cf7-ab8a-f4eec8a16e33&banner_id=1" \
		 -H "X-Tenant-Name: admin" \
		 -H "X-Auth-Token: $token" 
}

function upload-slogan() {
	curl -X POST "http://localhost:5000/api/v1/slogan" \
		 -H "X-Tenant-Name: admin" \
		 -H "X-Auth-Token: $token" \
		 -d '{
		 	"moodle_id": "d7ba7a3c-cb4e-4cf7-ab8a-f4eec8a16e33",
			"slogan_id": "3",
		 	"slogan": "hah"
		 }'
}

function upload-course() {
	curl -X POST "http://lms-service:5000/api/v1/pre-installed-course" \
		 -H "X-Tenant-Name: admin" \
		 -H "X-Auth-Token: $token" \
		 -d '{
		 	"id": "ffd90bce-5d60-4130-8c4b-b20b584cf2c1",
			"pre_installed_course": [1,2,3,4,5,6,7,8]
		 }'
}

function update-sitename() {
	curl -X POST "http://localhost:5000/api/v1/sitename" \
		 -H "X-Tenant-Name: admin" \
		 -H "X-Auth-Token: $token" \
		 -d '{
		 	"moodle_id": "2796701e-27c3-45d5-9746-aec55fc82360",
			"site_name": "gg"
		 }'
}



$*
