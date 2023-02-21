package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"moodle/config"
	"net/http"
)

// sendmail by call to external api
// this email endpoint sending SendMail
// localhost:5000/v1/lms/create_success
// this is payload
// {
//     "account": {
//         "mail": "{{mail}}"
//     },
//     "website_name": "duyseeexy",
//     "ip_address": "1.1.1.1",
//     "package": "100 CCU/ 10GB/ 5 Backups",
//     "username": "duysexy",
//     "password": "password1",
//     "cc": [],
//     "bcc": []
// }


type SendmailPayload struct {
	Account     Account     `json:"account"`
	WebsiteName string      `json:"website_name"`
	IPAddress   string      `json:"ip_address"`
	Package     string      `json:"package"`
	Username    string      `json:"username"`
	Password    string      `json:"password"`
	CC          []string    `json:"cc"`
	BCC         []string    `json:"bcc"`
}

type Account struct {
	Mail string `json:"mail"`
}

func SendMail(mail, websiteName, ipAddress, packageInfo, username, password string) error {
	url := "https://" + config.SENDMAILDOMAIN + ":5000/v1/lms/create_success"
	payload := SendmailPayload{
		Account: Account{
			Mail: mail,
		},
		WebsiteName: websiteName,
		IPAddress:   ipAddress,
		Package:     packageInfo,
		Username:    username,
		Password:    password,
		CC:          []string{},
		BCC:         []string{"duynn@bizflycloud.vn"},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("send mail failed")
	}
	return nil
}


