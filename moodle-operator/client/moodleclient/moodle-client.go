// zaddok moodle

package main

import (
	"github.com/zaddok/moodle"
	"fmt"
)


func main() {
	token := "21554247a241aa407a456bcd9c8be1eb"
	l := moodle.NewLogin(token)
	api := moodle.NewMoodleApi("http://moodle.example.com/webservice/rest/server.php", token, l)
	api.GetCourses()
}
