#### Create Moodle
<details>
<summary><code>POST</code> <code>(/api/v1/moodles)</code>Create Moodle</summary> 

##### Header
> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | X-Tenant-Name		   |  required |		string				|	email																|
> | X-Auth-Token		   |  required |		string				|	token																|

##### Body

> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | email	 |  required   |   string  |			N/A				|																		|
> | website_name		   |  required |			string			| N/A																	|
> | pre_installed_course   |  required |			[]int			| N/A																	|
> | packages_name		   |  required |			string			| N/A																	|
> | autoscale			   |  required |			bool			| N/A																	|
> | documents_storage_extra|  required |			int				| N/A																	|


##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `201`         | `application/json`				  | <pre lang="json">{<br>    "id":"49404390-ac3a-47f0-83ea-4820bb29a268",<br>    "email":"duynn@bizflycloud.vn",<br>    "ip":"Provisioning",<br>    "name":"test002",<br>    "lb_name":"kube_service_6o0cn9lv42livqek_49404390-ac3a-47f0-83ea-4820bb29a268_moodle-service",<br>    "website_name":"test002.lms.bizflycloud.vn",<br>    "pre_installed_course":[1, 2],<br>    "packages":{"name":"100CCU",<br>    "ccu":100,<br>    "account_max":4000,<br>    "document_storage":50,<br>    "moodle_version":"4.0.1",<br>    "backup_num":4,<br>    "ccu_extra_max":400,<br>    "document_storage_extra_max":2000},<br>    "autoscale":false,<br>    "documents_storage_extra":100,<br>    "status":"Creating",<br>    "created_at":"2023-01-11T10:14:59.406757545+07:00",<br>    "updated_at":"2023-01-11T10:14:59.406757751+07:00"<br>}</pre>|
> | `400`         | `application/json`                |<pre lang="json">{<br>    "error": Message<br>}</pre>											|

##### Example cURL

```curl
curl -i -X POST "http://localhost:5000/api/v1/moodles" \
		 -H "Content-Type: application/json" \
         -H "X-Tenant-Name: duynn@bizflycloud.vn" \
		 -H "X-Auth-Token: $TOKEN" \
		 -d '{
			 "email": "duydeptrai@bizflycloud.vn",
			 "website_name": test001",
			 "pre_installed_course": [1,2],
			 "packages_name": "100CCU",
			 "autoscale": true,
			 "documents_storage_extra": 100
		 }'

```

</details>

------------------------------------------------------------------------------------------

#### Listing Moodle, Get Moodle, Search Moodle, List Courses

<details>
<summary><code>GET</code> <code>(/api/v1/moodles)</code>List Moodle</summary> 

##### Header
> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | X-Tenant-Name		   |  required |		string				|	email																|
> | X-Auth-Token		   |  required |		string				|	token																|

##### Parameters
> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | page				   |  required |			int				| N/A																	|
> | limit				   |  required |			int				| N/A																	|
> | email	 |  required   |  string	| N/A						|																		|

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `application/json`				  | <pre lang="json">{<br>    "total":2,<br>    "pages":1,<br>    "page":1,<br>    "limit":10,<br>    "moodles":[{"id":"366561ed-de0b-422c-ac00-6556881cecab",<br>    "email":"duynn@bizflycloud.vn",<br>    "ip":"14.225.27.21",<br>    "name":"test001",<br>    "lb_name":"kube_service_6o0cn9lv42livqek_366561ed-de0b-422c-ac00-6556881cecab_moodle-service",<br>    "website_name":"test001.lms.bizflycloud.vn",<br>    "pre_installed_course":[1, 2],<br>    "packages":{"name":"100CCU",<br>    "ccu":100,<br>    "account_max":4000,<br>    "document_storage":50,<br>    "moodle_version":"4.0.1",<br>    "backup_num":4,<br>    "ccu_extra_max":400,<br>    "document_storage_extra_max":2000},<br>    "autoscale":false,<br>    "documents_storage_extra":100,<br>    "status":"ONLINE",<br>    "created_at":"2023-01-09T02:25:07.463Z",<br>    "updated_at":"2023-01-09T04:51:15.387Z"},<br>    {"id":"49404390-ac3a-47f0-83ea-4820bb29a268",<br>    "email":"duynn@bizflycloud.vn",<br>    "ip":"14.225.27.21",<br>    "name":"test002",<br>    "lb_name":"kube_service_6o0cn9lv42livqek_49404390-ac3a-47f0-83ea-4820bb29a268_moodle-service",<br>    "website_name":"test002.lms.bizflycloud.vn",<br>    "pre_installed_course":[1, 2],<br>    "packages":{"name":"100CCU",<br>    "ccu":100,<br>    "account_max":4000,<br>    "document_storage":50,<br>    "moodle_version":"4.0.1",<br>    "backup_num":4,<br>    "ccu_extra_max":400,<br>    "document_storage_extra_max":2000},<br>    "autoscale":false,<br>    "documents_storage_extra":100,<br>    "status":"ONLINE",<br>    "created_at":"2023-01-11T03:14:59.406Z",<br>    "updated_at":"2023-01-11T03:48:57.885Z"}]}</pre>|
> | `400`         | `application/json`                | <pre lang="json">{<br>    "error": Message<br>}</pre>				|


##### Example cURL

```curl
curl -iX GET "http://localhost:5000/api/v1/moodles?page=1&limit=10" \
	 -H "X-Tenant-Name: duynn@bizflycloud.vn" \
	 -H "X-Auth-Token: $TOKEN"
```

</details>

<details>
<summary><code>GET</code> <code>(/api/v1/moodles)</code>Get Moodle By Id</summary>

##### Header
> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | X-Tenant-Name		   |  required |		string				|	email																|
> | X-Auth-Token		   |  required |		string				|	token																|

##### Parameters
> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | id					   |  required |			string			| N/A																	|

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `application/json`				  | <pre lang="json">{<br>    "id":"49404390-ac3a-47f0-83ea-4820bb29a268",<br>    "email":"duynn@bizflycloud.vn",<br>    "ip":"Provisioning",<br>    "name":"test002",<br>    "lb_name":"kube_service_6o0cn9lv42livqek_49404390-ac3a-47f0-83ea-4820bb29a268_moodle-service",<br>    "website_name":"test002.lms.bizflycloud.vn",<br>    "pre_installed_course":[1, 2],<br>    "packages":{"name":"100CCU",<br>    "ccu":100,<br>    "account_max":4000,<br>    "document_storage":50,<br>    "moodle_version":"4.0.1",<br>    "backup_num":4,<br>    "ccu_extra_max":400,<br>    "document_storage_extra_max":2000},<br>    "autoscale":false,<br>    "documents_storage_extra":100,<br>    "status":"Creating",<br>    "created_at":"2023-01-11T10:14:59.406757545+07:00",<br>    "updated_at":"2023-01-11T10:14:59.406757751+07:00"}</pre>|
> | `400`         | `application/json`                | <pre lang="json">{<br>    "error": Message<br>}</pre>				|


##### Example cURL

```curl
curl -iX GET "http://localhost:5000/api/v1/moodles?id=366561ed-de0b-422c-ac00-6556881cecab" \
	 -H "X-Tenant-Name: duynn@bizflycloud.vn" \
	 -H "X-Auth-Token: $TOKEN"
```

</details>
<details>
<summary><code>GET</code> <code>(/api/v1/moodles)</code>Search Moodle By Name</summary>

##### Header
> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | X-Tenant-Name		   |  required |		string				|	email																|
> | X-Auth-Token		   |  required |		string				|	token																|

##### Parameters
> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | email				   |  required |			string			| N/A																	|
> | search				   |  required |			string			| N/A																	|
> | page				   |  required |			int				| N/A																	|
> | limit				   |  required |			int				| N/A																	|

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `application/json`				  | <pre lang="json">{<br>    "total":2,<br>    "pages":1,<br>    "page":1,<br>    "limit":10,<br>    "moodles":[{"id":"366561ed-de0b-422c-ac00-6556881cecab",<br>    "email":"duynn@bizflycloud.vn",<br>    "ip":"14.225.27.21",<br>    "name":"test001",<br>    "lb_name":"kube_service_6o0cn9lv42livqek_366561ed-de0b-422c-ac00-6556881cecab_moodle-service",<br>    "website_name":"test001.lms.bizflycloud.vn",<br>    "pre_installed_course":[1, 2],<br>    "packages":{"name":"100CCU",<br>    "ccu":100,<br>    "account_max":4000,<br>    "document_storage":50,<br>    "moodle_version":"4.0.1",<br>    "backup_num":4,<br>    "ccu_extra_max":400,<br>    "document_storage_extra_max":2000},<br>    "autoscale":false,<br>    "documents_storage_extra":100,<br>    "status":"ONLINE",<br>    "created_at":"2023-01-09T02:25:07.463Z",<br>    "updated_at":"2023-01-09T04:51:15.387Z"},<br>    {"id":"49404390-ac3a-47f0-83ea-4820bb29a268",<br>    "email":"duynn@bizflycloud.vn",<br>    "ip":"14.225.27.21",<br>    "name":"test002",<br>    "lb_name":"kube_service_6o0cn9lv42livqek_49404390-ac3a-47f0-83ea-4820bb29a268_moodle-service",<br>    "website_name":"test002.lms.bizflycloud.vn",<br>    "pre_installed_course":[1, 2],<br>    "packages":{"name":"100CCU",<br>    "ccu":100,<br>    "account_max":4000,<br>    "document_storage":50,<br>    "moodle_version":"4.0.1",<br>    "backup_num":4,<br>    "ccu_extra_max":400,<br>    "document_storage_extra_max":2000},<br>    "autoscale":false,<br>    "documents_storage_extra":100,<br>    "status":"ONLINE",<br>    "created_at":"2023-01-11T03:14:59.406Z",<br>    "updated_at":"2023-01-11T03:48:57.885Z"}]}</pre>|
> | `400`         | `application/json`                | <pre lang="json">{<br>    "error": Message<br>}</pre>				|


##### Example cURL

```curl
curl -iX GET "http://localhost:5000/api/v1/moodles?search=test&page=1&limit=10" \
	 -H "X-Tenant-Name: duynn@bizflycloud.vn" \
	 -H "X-Auth-Token: $TOKEN"
```

</details>



<details>
<summary><code>GET</code> <code>(/api/v1/courses)</code>List Moodle Course</summary>

##### Header
> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | X-Tenant-Name		   |  required |		string				|	email																|
> | X-Auth-Token		   |  required |		string				|	token																|
##### Parameters
> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | page				   |  required |			int				| N/A																	|
> | limit				   |  required |			int				| N/A																	|

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `application/json`				  | <pre lang="json">{"total":8<br>    "pages":1<br>    "page":1<br>    "limit":10<br>    "courses":[{"id":1<br>    "name":"Hội nhập"<br>    "content":"Khóa học giúp nhân viên nhanh chóng nắm bắt được các thông tin chung của doanh nghiệp"}<br>    {"id":2<br>    "name":"Định hướng"<br>    "content":"Giúp doanh nghiệp có được cái nhìn tổng quát về năng lực thực tế của nhân viên"}<br>    {"id":3<br>    "name":"Phát triển chuyên môn"<br>    "content":"Giúp doanh nghiệp có được cái nhìn tổng quát về năng lực thực tế của nhân viên"}<br>    {"id":4<br>    "name":"Phát triển kỹ năng mềm"<br>    "content":"Hình thức đào tạo nhân sự này thường được áp dụng tại các doanh nghiệp"}<br>    {"id":5<br>    "name":"Dịch vụ \u0026 sản phẩm"<br>    "content":"Khóa học giúp nhân viên nhanh chóng nắm bắt được các thông tin chung của doanh nghiệp"}<br>    {"id":6<br>    "name":"Tiêu chuẩn và chất lượng"<br>    "content":"Giúp doanh nghiệp có được cái nhìn tổng quát về năng lực thực tế của nhân viên"}<br>    {"id":7<br>    "name":"An toàn và lao động"<br>    "content":"Giúp doanh nghiệp có được cái nhìn tổng quát về năng lực thực tế của nhân viên"}<br>    {"id":8<br>    "name":"Đội nhóm"<br>    "content":"Hình thức đào tạo nhân sự này thường được áp dụng tại các doanh nghiệp"}]}</pre>|
> | `400`         | `application/json`                | <pre lang="json">{<br>    "error": Message<br>}</pre>				|


##### Example cURL

```curl
curl -H "X-Tenant-Name: duynn@bizflycloud.vn" \
	 -H "X-Auth-Token: $token" \
	 -iX GET "http://localhost:5000/api/v1/courses?page=1&limit=10"
```

</details>
------------------------------------------------------------------------------------------

#### Update Packages,  Update Pre-installed Course,  Update AutoScale, Update Storage Extra
<details>
<summary><code>PUT</code> <code>(/api/v1/packages)</code>Update Packages</summary>

##### Header
> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | X-Tenant-Name		   |  required |		string				|	email																|
> | X-Auth-Token		   |  required |		string				|	token																|
##### Body
> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | moodle_id			   |  required |			string			| moodle_id																|
> | packages_name		   |  required |			int				| 100CCU, 200CCU, 300CCU												|

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `201`         | `application/json`				  | <pre lang="json">{"id":"366561ed-de0b-422c-ac00-6556881cecab"<br>    "email":"duynn@bizflycloud.vn"<br>    "ip":"14.225.27.21"<br>    "name":"test001"<br>    "lb_name":"kube_service_6o0cn9lv42livqek_366561ed-de0b-422c-ac00-6556881cecab_moodle-service"<br>    "website_name":"test001.lms.bizflycloud.vn"<br>    "pre_installed_course":[1<br>    2]<br>    "packages":{"name":"200CCU"<br>    "ccu":200<br>    "account_max":8000<br>    "document_storage":100<br>    "moodle_version":"4.0.1"<br>    "backup_num":4<br>    "ccu_extra_max":800<br>    "document_storage_extra_max":4000}<br>    "autoscale":false<br>    "documents_storage_extra":100<br>    "status":"ONLINE"<br>    "created_at":"2023-01-09T02:25:07.463Z"<br>    "updated_at":"2023-01-09T04:51:15.387Z"}</pre>	  |
> | `400`         | `application/json`                | <pre lang="json">{<br>    "error": Message<br>}</pre>				|



##### Example cURL

```curl
curl -i -X PUT "http://localhost:5000/api/v1/packages" \
	 -H "X-Auth-Token: $token" \
	 -H "X-Tenant-Name: duynn@bizflycloud.vn" \
	 -H "Content-Type: application/json" \
	 -d "{
		\"moodle_id\": \"$1\",
		\"packages_name\": \"$2\"
	}"
```

</details>

<details>
<summary><code>PUT</code> <code>(/api/v1/pre-installed-course)</code>Update Pre-installed Course</summary>

##### Header
> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | X-Tenant-Name		   |  required |		string				|	email																|
> | X-Auth-Token		   |  required |		string				|	token																|
##### Body
> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | moodle_id			   |  required |			string			| moodle_id																|
> | pre_installed_course   |  required |			[]int			| [1,2,3]

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `201`         | `application/json`				  | <pre lang="json">{"id":"366561ed-de0b-422c-ac00-6556881cecab"<br>    "email":"duynn@bizflycloud.vn"<br>    "ip":"14.225.27.21"<br>    "name":"test001"<br>    "lb_name":"kube_service_6o0cn9lv42livqek_366561ed-de0b-422c-ac00-6556881cecab_moodle-service"<br>    "website_name":"test001.lms.bizflycloud.vn"<br>    "pre_installed_course":[1<br>    2]<br>    "packages":{"name":"200CCU"<br>    "ccu":200<br>    "account_max":8000<br>    "document_storage":100<br>    "moodle_version":"4.0.1"<br>    "backup_num":4<br>    "ccu_extra_max":800<br>    "document_storage_extra_max":4000}<br>    "autoscale":false<br>    "documents_storage_extra":100<br>    "status":"ONLINE"<br>    "created_at":"2023-01-09T02:25:07.463Z"<br>    "updated_at":"2023-01-09T04:51:15.387Z"}</pre>	  |
> | `400`         | `application/json`                | <pre lang="json">{<br>    "error": Message<br>}</pre>				|




##### Example cURL

```curl
curl -i -X PUT "http://localhost:5000/api/v1/pre-installed-course" \
	 -H "X-Auth-Token: $token" \
	 -H "X-Tenant-Name: duynn@bizflycloud.vn" \
	 -H "Content-Type: application/json" \
	 -d "{
		\"moodle_id\": \"$1\",
		\"pre_installed_course\": $2
	 }"
```

</details>



<details>
<summary><code>PUT</code> <code>(/api/v1/autoscale)</code>Update AutoScale</summary>

##### Header
> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | X-Tenant-Name		   |  required |		string				|	email																|
> | X-Auth-Token		   |  required |		string				|	token																|
##### Body
> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | moodle_id			   |  required |			string			| moodle_id																|
> | pre_installed_course   |  required |			[]int			| [1,2,3]

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `201`         | `application/json`				  | <pre lang="json">{"id":"366561ed-de0b-422c-ac00-6556881cecab"<br>    "email":"duynn@bizflycloud.vn"<br>    "ip":"14.225.27.21"<br>    "name":"test001"<br>    "lb_name":"kube_service_6o0cn9lv42livqek_366561ed-de0b-422c-ac00-6556881cecab_moodle-service"<br>    "website_name":"test001.lms.bizflycloud.vn"<br>    "pre_installed_course":[2<br>    3<br>    4]<br>    "packages":{"name":"200CCU"<br>    "ccu":200<br>    "account_max":8000<br>    "document_storage":100<br>    "moodle_version":"4.0.1"<br>    "backup_num":4<br>    "ccu_extra_max":800<br>    "document_storage_extra_max":4000}<br>    "autoscale":true<br>    "documents_storage_extra":100<br>    "status":"ONLINE"<br>    "created_at":"2023-01-09T02:25:07.463Z"<br>    "updated_at":"2023-01-09T04:51:15.387Z"}</pre>|
> | `400`         | `application/json`                | <pre lang="json">{<br>    "error": Message<br>}</pre>				  |



##### Example cURL

```curl
curl -i -X PUT "http://localhost:5000/api/v1/autoscale" \
	 -H "X-Auth-Token: $token" \
	 -H "X-Tenant-Name: duynn@bizflycloud.vn" \
	 -H "Content-Type: application/json" \
	 -d "{
		\"moodle_id\": \"$1\",
		\"autoscale\": $2
	 }"
```

</details>
<details>
<summary><code>PUT</code> <code>(/api/v1/document-storage-extra)</code>Update Storage Extra</summary>

##### Header
> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | X-Tenant-Name		   |  required |		string				|	email																|
> | X-Auth-Token		   |  required |		string				|	token																|
##### Body
> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | moodle_id			   |  required |			string			| moodle_id																|
> | documents_storage_extra|  required |			int				| 

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `201`         | `application/json`				  | <pre lang="json">{"id":"366561ed-de0b-422c-ac00-6556881cecab"<br>    "email":"duynn@bizflycloud.vn"<br>    "ip":"14.225.27.21"<br>    "name":"test001"<br>    "lb_name":"kube_service_6o0cn9lv42livqek_366561ed-de0b-422c-ac00-6556881cecab_moodle-service"<br>    "website_name":"test001.lms.bizflycloud.vn"<br>    "pre_installed_course":[2<br>    3<br>    4]<br>    "packages":{"name":"200CCU"<br>    "ccu":200<br>    "account_max":8000<br>    "document_storage":100<br>    "moodle_version":"4.0.1"<br>    "backup_num":4<br>    "ccu_extra_max":800<br>    "document_storage_extra_max":4000}<br>    "autoscale":true<br>    "documents_storage_extra":500<br>    "status":"ONLINE"<br>    "created_at":"2023-01-09T02:25:07.463Z"<br>    "updated_at":"2023-01-09T04:51:15.387Z"}</pre>|
> | `400`         | `application/json`                | <pre lang="json">{<br>    "error": Message<br>}</pre>				  |



##### Example cURL

```curl
curl -i -X PUT "http://localhost:5000/api/v1/document-storage-extra" \
	 -H "Content-Type: application/json" \
	 -H "X-Auth-Token: $token" \
	 -H "X-Tenant-Name: duynn@bizflycloud.vn" \
	 -d "{
		\"moodle_id\": \"$1\",
		\"documents_storage_extra\": $2
	 }"
```

</details>

------------------------------------------------------------------------------------------
#### Delete Moodle
<details>
<summary><code>DELETE</code> <code>(/api/v1/moodles)</code>Delete Moodle</summary>

##### Header
> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | X-Tenant-Name		   |  required |		string				|	email																|
> | X-Auth-Token		   |  required |		string				|	token																|
##### Body
> | name				   |  type     | data type					| description                                                           |
> |------------------------|-----------|----------------------------|-----------------------------------------------------------------------|
> | moodle_id			   |  required |			string			| moodle_id																|

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `204`         | `application/json`				  |																		|
> | `400`         | `application/json`                | <pre lang="json">{<br>    "error": Message<br>}</pre>				|



##### Example cURL

```curl
curl -i -X DELETE "http://localhost:5000/api/v1/moodles" \
	 -H "Content-Type: application/json" \
	 -H "X-Auth-Token: $token" \
	 -H "X-Tenant-Name: duynn@bizflycloud.vn" \
	 -d "{
		\"moodle_id\": \"$1\"
	 }"

```

</details>
