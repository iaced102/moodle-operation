package client

import (
	"fmt"
	"io"
	"log"
	"moodle/config"
	"net/http"
	"strings"
)

func SendMailCreate(mail, websiteName, ipAddress, packageInfo, username, password, createdat string) error {
	client := &http.Client{}
	var data = strings.NewReader(`{
			"account": {
				"mail": "` + mail + `"
			},
			"website_name": "` + websiteName + `",
			"ip_address": "` + ipAddress + `",
			"package": "` + packageInfo + `",
			"username": "` + username + `",
			"password": "` + password + `",
			"cc": [],
			"created_at": "` + createdat + `",
			"bcc": []
		}`)
	req, err := http.NewRequest("POST", config.SENDMAILDOMAIN + "/v1/lms/create_success", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// if status code is not 200, return body as error
	if resp.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf(string(bodyBytes))
	}
	return nil
}


func SendMailDelete(mail, websiteName, ipAddress, packageInfo, createdat string) error {
	client := &http.Client{}
	var data = strings.NewReader(`{
			"account": {
				"mail": "` + mail + `"
			},
			"website_name": "` + websiteName + `",
			"ip_address": "` + ipAddress + `",
			"package": "` + packageInfo + `",
			"cc": [],
			"deleted_at": "` + createdat + `",
			"bcc": []
		}`)
	req, err := http.NewRequest("POST", config.SENDMAILDOMAIN + "/v1/lms/delete_success", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// if status code is not 200, return body as error
	if resp.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf(string(bodyBytes))
	}
	return nil
}
