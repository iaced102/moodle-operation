#!/bin/bash

# list moodle
function moodle-list() {
	curl -iX GET "http://localhost:5000/api/moodles?userid=$1" 
}

# create moodle
function moodle-create() {
	curl -i -X POST "http://localhost:5000/api/moodles" \
	-H "Content-Type: application/json" \
	-d "{
		\"name\": \"test\",
		\"ccu\": 100,
		\"userid\": \"$1\",
		\"theme\": \"$2\"
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
			\"id\": \"$1\"
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
