package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"moodle/config"
	"net/http"
)

func NewLBClient() *LBClient {
	return &LBClient{}
}

type LBClient struct {
}



type LoadBalancer struct {
	ID                 string       `json:"id"`
	FlavorID           string       `json:"flavor_id"`
	Description        string       `json:"description"`
	Provider           string       `json:"provider"`
	UpdatedAt          string       `json:"updated_at"`
	Listeners          []resourceID `json:"listeners"`
	VipSubnetID        string       `json:"vip_subnet_id"`
	ProjectID          string       `json:"project_id"`
	VipQosPolicyID     string       `json:"vip_qos_policy_id"`
	VipNetworkID       string       `json:"vip_network_id"`
	NetworkType        string       `json:"network_type"`
	VipAddress         string       `json:"vip_address"`
	VipPortID          string       `json:"vip_port_id"`
	AdminStateUp       bool         `json:"admin_state_up"`
	Name               string       `json:"name"`
	OperatingStatus    string       `json:"operating_status"`
	ProvisioningStatus string       `json:"provisioning_status"`
	Pools              []resourceID `json:"pools"`
	Type               string       `json:"type"`
	TenantID           string       `json:"tenant_id"`
	CreatedAt          string       `json:"created_at"`
}

type resourceID struct {
	ID string
}


// return list of instances
func (m *LBClient) GetLBInstances(token string) ([]*LoadBalancer, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://hn.manage.bizflycloud.vn/api/loadbalancers/loadbalancers", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Region-Name", "HaNoi")
	req.Header.Set("X-Tenant-Name", config.USERNAME)
	req.Header.Set("X-Auth-Token", token)
	resp, err := client.Do(req)
	if resp.StatusCode == 401 {
		log.Fatal("Unauthorized")
	}
	// log.Println(resp)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var data struct {
		LB []*LoadBalancer `json:"loadbalancers"`
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bodyText, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data.LB, nil
}


