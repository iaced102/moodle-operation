#!/bin/bash

# list moodle
function moodle-list() {
	curl -iX GET "http://localhost:5000/api/moodles?user_id=$1" 
}

# create moodle
function moodle-create() {
	echo $1 $2 $3 $4 $5 $6
	curl -i -X POST "http://localhost:5000/api/moodles" \
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

# change packages moodle
function moodle-change-packages() {
	curl -i -X PUT "http://localhost:5000/api/moodles/packages" \
	-H "Content-Type: application/json" \
	-d "{
		\"id\": \"$1\",
		\"packages_name\": \"$2\"
	}"
}


# change autoscale moodle
function moodle-change-autoscale() {
	curl -i -X PUT "http://localhost:5000/api/moodles/autoscale" \
	-H "Content-Type: application/json" \
	-d "{
		\"id\": \"$1\",
		\"autoscale\": $2
	}"
}

# change documents_storage_extra moodle
function moodle-change-document-storage-extra() {
	curl -i -X PUT "http://localhost:5000/api/moodles/document-storage-extra" \
	-H "Content-Type: application/json" \
	-d "{
		\"id\": \"$1\",
		\"documents_storage_extra\": $2
	}"
}

# change pre_installed_course moodle
function moodle-change-pre-installed-course() {
	curl -i -X PUT "http://localhost:5000/api/moodles/pre-installed-course" \
	-H "Content-Type: application/json" \
	-d "{
		\"id\": \"$1\",
		\"pre_installed_course\": $2
	}"
}

# delete moodle
function moodle-delete() {
	curl -i -X DELETE "http://localhost:5000/api/moodles" \
		-H "Content-Type: application/json" \
		-d "{
			\"id\": \"$1\"
		}"
}

# scale moodle
function moodle-scale() {
	curl -i -X PUT "http://localhost:5000/api/moodles" \
		-H "Content-Type: application/json" \
		-d "{
			\"id\": \"$1\",
			\"replicas\": $2
		}"
}

# create user
function user-add() {
	curl -i -X POST http://localhost:5000/api/users \
		-H "Content-Type: application/json" \
		-d "{
			\"email\": \"$1\"
		}"
}

# delete user
function user-delete() {
	curl -i -X DELETE "http://localhost:5000/api/users" \
		-H "Content-Type: application/json" \
		-d "{
			\"id\": \"$1\"
		}"
}

$*
