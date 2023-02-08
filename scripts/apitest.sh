#!/bin/bash
token="gAAAAABjvhxZwCj1n8fFu-nFALaMr4KKZ7G3UPI7pAETzVxGa4pYuWdUelC6DcjLSRd8rHXIWVsya4q0qijXmOiwPbJzfNOSxBthvRZmGBxfc46rHefbWpwXq-EqFAsxO5MDFPYjUsBocTwCtb6UFdP8ruEdAWxWQ9YQl2KPAuyP_bprvmZwUkI"

# create moodle
function moodle-create() {
	curl -s -i -X POST "http://localhost:5000/api/v1/moodles" \
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
	curl -s -H "X-Tenant-Name: admin" \
		 -H "X-Auth-Token: $token" \
		 -iX GET "http://localhost:5000/api/v1/moodles?page=1&limit=10&email=$1"
}

# get moodle
function moodle-get() {
	curl -s -H "X-Tenant-Name: admin" \
		 -H "X-Auth-Token: $token" \
		 -iX GET "http://localhost:5000/api/v1/moodles?id=$1" 
}

# search moodle
function search() {
	curl -s -X GET "http://localhost:5000/api/v1/moodles?search=$1&page=1&limit=10&email=duynn@bizflycloud.vn" \
		 -H "X-Auth-Token: $token" \
		 -H "X-Tenant-Name: admin"
}

# change packages moodle
function moodle-change-packages() {
	curl -s -i -X PUT "http://localhost:5000/api/v1/packages" \
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
	curl -s -i -X PUT "http://localhost:5000/api/v1/autoscale" \
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
	curl -s -i -X PUT "http://localhost:5000/api/v1/document-storage-extra" \
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
	curl -s -i -X PUT "http://localhost:5000/api/v1/pre-installed-course" \
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
	curl -s -i -X DELETE "http://localhost:5000/api/v1/moodles" \
		 -H "Content-Type: application/json" \
		 -H "X-Auth-Token: $token" \
		 -H "X-Tenant-Name: admin" \
		 -d "{
		 	\"moodle_id\": \"$1\"
		 }"
}

# list course
function course-list() {
	curl -s -H "X-Tenant-Name: admin" \
		 -H "X-Auth-Token: $token" \
		 -iX GET "http://localhost:5000/api/v1/courses?page=1&limit=10"
}
# list pacakge
function package-list() {
	curl -s -H "X-Tenant-Name: admin" \
		 -H "X-Auth-Token: $token" \
		 -iX GET "http://localhost:5000/api/v1/packages?page=1&limit=10"
}

# upload file
function upload-logo() {
	curl -F "file=@/home/duy/Pictures/duydeptrai.jpg" "http://localhost:5000/api/v1/logo?moodle_id=testid" \
		 -H "X-Tenant-Name: admin" \
		 -H "X-Auth-Token: $token" 
}

$*
