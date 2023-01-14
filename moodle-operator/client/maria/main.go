package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"moodle/config"
	"net/http"
	"strings"
)


type MariaClient struct {
}

func NewMariaClient() *MariaClient {
	return &MariaClient{}
}

type Session struct {
	expire_at string
	token     string
}

func (m *MariaClient) GetSession() *Session {
	client := &http.Client{}
	var data = strings.NewReader(`{
            "username": "` + config.USERNAME+ `",
            "password": "` + config.PASSWORD+ `",
            "auth_method": "password"
        }`)
	req, err := http.NewRequest("POST", "https://manage.bizflycloud.vn/api/token", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var session Session
	session.token = string(bodyText)
	return &session
}

type CloudDatabaseInstance struct {
	AutoScaling  CloudDatabaseAutoScaling `json:"autoscaling"`
	CreatedAt    string                   `json:"created"`
	Datastore    CloudDatabaseDatastore   `json:"datastore"`
	Description  string                   `json:"description"`
	DNS          CloudDatabaseDNS         `json:"dns"`
	ID           string                   `json:"id"`
	Logs         CloudDatabaseLog         `json:"logs"`
	Message      string                   `json:"message"`
	Name         string                   `json:"name"`
	Networks     []CloudDatabaseNetworks  `json:"networks"`
	Nodes        []CloudDatabaseNode      `json:"nodes"`
	ProjectID    string                   `json:"project_id"`
	PublicAccess bool                     `json:"public_access"`
	Status       string                   `json:"status"`
	TaskID       string                   `json:"task_id"`
	Volume       CloudDatabaseVolume      `json:"volume"`
}

type CloudDatabaseDatastore struct {
	Type      string `json:"type,omitempty"`
	Name      string `json:"name,omitempty"`
	ID        string `json:"id,omitempty"`
	VersionID string `json:"version_id,omitempty"`
}

type CloudDatabaseDNS struct {
	Private string `json:"private"`
	Public  string `json:"public"`
	SRV     string `json:"srv"`
}

type CloudDatabaseVolume struct {
	Size int     `json:"size"`
	Used float32 `json:"used"`
}

type CloudDatabaseAddressesDetail struct {
	IPAddress   string `json:"ip_address"`
	NetworkName string `json:"network_name"`
}

type CloudDatabaseAddresses struct {
	Private []CloudDatabaseAddressesDetail `json:"private"`
	Public  []CloudDatabaseAddressesDetail `json:"public"`
}

type CloudDatabaseNode struct {
	Addresses        CloudDatabaseAddresses `json:"addresses"`
	AvailabilityZone string                 `json:"availability_zone"`
	CreatedAt        string                 `json:"created_at"`
	Datastore        CloudDatabaseDatastore `json:"datastore"`
	Description      string                 `json:"description"`
	DNS              CloudDatabaseDNS       `json:"dns"`
	Flavor           string                 `json:"flavor.id"`
	ID               string                 `json:"id"`
	InstanceID       string                 `json:"instance_id"`
	Message          string                 `json:"message"`
	Name             string                 `json:"name"`
	RegionName       string                 `json:"region_name"`
	Replicas         []CloudDatabaseNode    `json:"replicas"`
	Role             string                 `json:"role"`
	Status           string                 `json:"status"`
	TaskID           string                 `json:"task_id"`
	Volume           CloudDatabaseVolume    `json:"volume"`
}

type CloudDatabaseLog struct {
	Enable string `json:"enable"`
	Name   string `json:"name"`
}

type CloudDatabaseAutoScalingAlarms struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ReceiverID string `json:"receiver_id"`
}

type CloudDatabaseAutoScalingReceivers struct {
	Action string `json:"action"`
	ID     string `json:"id"`
	Name   string `json:"name"`
}

type CloudDatabaseAutoScalingVolume struct {
	Limited   int `json:"limited"`
	Threshold int `json:"threshold"`
}

type CloudDatabaseAutoScaling struct {
	Alarms    []CloudDatabaseAutoScalingAlarms    `json:"alarms,omitempty"`
	Enable    bool                                `json:"enable,omitempty"`
	Receivers []CloudDatabaseAutoScalingReceivers `json:"receivers,omitempty"`
	Volume    CloudDatabaseAutoScalingVolume      `json:"volume,omitempty"`
}

type CloudDatabaseNetworks struct {
	NetworkID string `json:"network_id"`
}


// return list of instances
func (m *MariaClient) GetMariaInstances() ([]*CloudDatabaseInstance , error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://hn.manage.bizflycloud.vn/api/cloud-database/instances", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Region-Name", "HaNoi")
	req.Header.Set("X-Tenant-Name", config.USERNAME)
	req.Header.Set("X-Auth-Token", config.BIZFLYCLOUD_TOKEN)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// return message as string if request fail 
	if resp.StatusCode != 200 {
		bodyText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		return nil, errors.New(string(bodyText))
	}
	var data struct {
		Instances []*CloudDatabaseInstance `json:"instances"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal(err)
	}
	// for _, instance := range data.Instances {
	// 	fmt.Println(instance)
	// }

	return data.Instances, nil
}


// create instance
func (m *MariaClient) CreateInstance(name string) error {
	client := &http.Client{}
	var data = strings.NewReader(`{"networks":[{"network_id":"7ef80d87-f2b0-4409-9b74-1b10e84fee1e"}],"public_access":true,"datastore":{"type":"MariaDB","version_id":"550aebf7-df97-49f1-bf24-7cd7b69fa365"},"autoscaling":{"enable":false,"volume":{"threshold":80,"limited":260}},"volume_size":40,"flavor_name":"2c_4g","availability_zone":"HN1","name":"` + name + `"}`)
	req, err := http.NewRequest("POST", "https://hn.manage.bizflycloud.vn/api/cloud-database/instances", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Region-Name", "HaNoi")
	req.Header.Set("X-Tenant-Name", config.USERNAME)
	req.Header.Set("X-Auth-Token", config.BIZFLYCLOUD_TOKEN)
	req.Header.Set("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)
	return nil
}


// delete instance
func (m *MariaClient) DeleteInstance(instanceID string) error {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", "https://hn.manage.bizflycloud.vn/api/cloud-database/instances/"+instanceID, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Region-Name", "HaNoi")
	req.Header.Set("X-Tenant-Name", config.USERNAME)
	req.Header.Set("X-Auth-Token", config.BIZFLYCLOUD_TOKEN)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)
	return nil
}

// create instance from backup
func (m *MariaClient) CreateInstanceFromBackup(name string, backupID string) error {
	client := &http.Client{}
	var data = strings.NewReader(`{"networks":[{"network_id":"7ef80d87-f2b0-4409-9b74-1b10e84fee1e"}],"public_access":false,"datastore":{"type":"MariaDB","version_id":"550aebf7-df97-49f1-bf24-7cd7b69fa365"},"autoscaling":{"enable":false,"volume":{"threshold":80,"limited":180}},"volume_size":40,"flavor_name":"1c_2g","availability_zone":"HN1","name":"` + name + `","backup_id":"` + backupID + `"}`)
	req, err := http.NewRequest("POST", "https://hn.manage.bizflycloud.vn/api/cloud-database/instances", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Region-Name", "HaNoi")
	req.Header.Set("X-Tenant-Name", config.USERNAME)
	req.Header.Set("X-Auth-Token", config.BIZFLYCLOUD_TOKEN)
	req.Header.Set("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// error if status code != 201
	if resp.StatusCode != 201 {
		return errors.New(string(bodyText))
	}
	return nil
}
