package hetznerrobot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type HetznerRobotClient struct {
	username string
	password string
}

func NewHetznerRobotClient(username string, password string) HetznerRobotClient {
	return HetznerRobotClient{
		username: username,
		password: password,
	}
}

type HetznerRobotFirewallResponse struct {
	Firewall HetznerRobotFirewall `json:"firewall"`
}

type HetznerRobotFirewall struct {
	IP                       string                    `json:"server_ip"`
	WhitelistHetznerServices bool                      `json:"whitelist_hos"`
	Status                   string                    `json:"status"`
	Rules                    HetznerRobotFirewallRules `json:"rules"`
}

type HetznerRobotFirewallRules struct {
	Input []HetznerRobotFirewallRule `json:"input"`
}

type HetznerRobotFirewallRule struct {
	Name     string `json:"name"`
	DstIp    string `json:"dst_ip"`
	DstPort  string `json:"dst_port"`
	SrcIp    string `json:"src_ip"`
	SrcPort  string `json:"src_port"`
	Protocol string `json:"protocol"`
	TCPFlags string `json:"tcp_flags"`
	Action   string `json:"action"`
}

func (c *HetznerRobotClient) getFirewall(ip string) (*HetznerRobotFirewall, error) {
	r, err := http.NewRequest("GET", fmt.Sprintf("https://robot-ws.your-server.de/firewall/%s", ip), nil)
	if err != nil {
		return nil, err
	}
	r.SetBasicAuth(c.username, c.password)

	client := http.Client{}
	response, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Hetzner response status %d\n%s", response.StatusCode, bytes)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Hetzner API response HTTP %d: %s", response.StatusCode, bytes)
	}
	firewall := HetznerRobotFirewallResponse{}
	if err = json.Unmarshal(bytes, &firewall); err != nil {
		return nil, err
	}
	return &firewall.Firewall, nil
}

func (c *HetznerRobotClient) setFirewall(firewall HetznerRobotFirewall) error {
	formParams := url.Values{}

	whitelistHOS := "false"
	if firewall.WhitelistHetznerServices {
		whitelistHOS = "true"
	}

	formParams.Set("whitelist_hos", whitelistHOS)
	formParams.Set("status", firewall.Status)

	for idx, rule := range firewall.Rules.Input {
		formParams.Set(fmt.Sprintf("rules[input][%d][%s]", idx, "ip_version"), "ipv4")
		formParams.Set(fmt.Sprintf("rules[input][%d][%s]", idx, "name"), rule.Name)
		formParams.Set(fmt.Sprintf("rules[input][%d][%s]", idx, "dst_ip"), rule.DstIp)
		formParams.Set(fmt.Sprintf("rules[input][%d][%s]", idx, "dst_port"), rule.DstPort)
		formParams.Set(fmt.Sprintf("rules[input][%d][%s]", idx, "src_ip"), rule.SrcIp)
		formParams.Set(fmt.Sprintf("rules[input][%d][%s]", idx, "src_port"), rule.SrcPort)
		if rule.Protocol != "" {
			formParams.Set(fmt.Sprintf("rules[input][%d][%s]", idx, "protocol"), rule.Protocol)
		}
		formParams.Set(fmt.Sprintf("rules[input][%d][%s]", idx, "action"), rule.Action)
	}

	encodedParams := formParams.Encode()
	log.Println(encodedParams)

	r, err := http.NewRequest("POST", fmt.Sprintf("https://robot-ws.your-server.de/firewall/%s", firewall.IP), strings.NewReader(encodedParams))
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.SetBasicAuth(c.username, c.password)

	client := http.Client{}
	response, err := client.Do(r)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Hetzner API response HTTP %d: %s", response.StatusCode, bytes)
	}

	return nil
}
