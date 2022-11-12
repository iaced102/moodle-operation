#!/bin/bash



token="gAAAAABjcaOJQxIBVlxL2vAnwiHyvsUpZikBBjfTJxJOBuAlOnb6JlRwoQ7XtTIYRHhm6ZlyEbWGhUl53c-diFI8XthmYvcJbE8-HEfA4Lbckc_96-vPvZpWrN8DTJEJRzktzMNOCjsESOdWpgQGAivys51Z5VJX3agPqcaJSkbbP8JOWzwHoDM"
username="duynn@bizflycloud.vn"
password="MzU2ZmE3ZDM5MmUwMDFkYThiMzkwMmE3"

get_token() {
    curl -si --location --request POST 'https://manage.bizflycloud.vn/api/token' \
        --header 'Content-Type: application/json' \
        -d "{
            \"username\": \"$username\",
            \"password\": \"$password\",
            \"auth_method\": \"password\"
        }" 
}

        # -o /dev/null -w '%{http_code}\n'

create_db() {
	curl 'https://hn.manage.bizflycloud.vn/api/cloud-database/instances' \
			-H "X-Region-Name: HaNoi" \
			-H "X-Tenant-Name: $username" \
			-H "X-Auth-Token: $token" \
			-H 'content-type: application/json' \
			--data-raw '{"networks":[{"network_id":"7ef80d87-f2b0-4409-9b74-1b10e84fee1e"}],"public_access":true,"datastore":{"type":"MariaDB","version_id":"550aebf7-df97-49f1-bf24-7cd7b69fa365"},"autoscaling":{"enable":false,"volume":{"threshold":80,"limited":260}},"volume_size":40,"flavor_name":"2c_4g","availability_zone":"HN1","name":"moodle009"}' \

	}

list_db() {
	curl 'https://hn.manage.bizflycloud.vn/api/cloud-database/instances?instance_name=&results_per_page=5&page=1' \
			-H "X-Region-Name: HaNoi" \
			-H "X-Tenant-Name: $username" \
			-H "X-Auth-Token: $token" 
	}

instance_id="aa9e9e97-4324-4f68-b6bf-55debff79d10"
delete_db() {
	curl "https://hn.manage.bizflycloud.vn/api/cloud-database/instances/$instance_id" \
		-X 'DELETE' \
		-H "X-Region-Name: HaNoi" \
		-H "X-Tenant-Name: $username" \
		-H "X-Auth-Token: $token" \
		--data-raw '{"purge_backup":true,"purge_autobackup":true}'
	}

$*



