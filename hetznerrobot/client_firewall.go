package hetznerrobot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

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
	DstIP    string `json:"dst_ip"`
	DstPort  string `json:"dst_port"`
	SrcIP    string `json:"src_ip"`
	SrcPort  string `json:"src_port"`
	Protocol string `json:"protocol"`
	TCPFlags string `json:"tcp_flags"`
	Action   string `json:"action"`
}

func (c *HetznerRobotClient) getFirewall(ip string) (*HetznerRobotFirewall, error) {

	bytes, err := c.makeAPICall("GET", fmt.Sprintf("%s/firewall/%s", c.url, ip), nil, http.StatusOK)
	if err != nil {
		return nil, err
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
		formParams.Set(fmt.Sprintf("rules[input][%d][%s]", idx, "dst_ip"), rule.DstIP)
		formParams.Set(fmt.Sprintf("rules[input][%d][%s]", idx, "dst_port"), rule.DstPort)
		formParams.Set(fmt.Sprintf("rules[input][%d][%s]", idx, "src_ip"), rule.SrcIP)
		formParams.Set(fmt.Sprintf("rules[input][%d][%s]", idx, "src_port"), rule.SrcPort)
		if rule.Protocol != "" {
			formParams.Set(fmt.Sprintf("rules[input][%d][%s]", idx, "protocol"), rule.Protocol)
		}
		formParams.Set(fmt.Sprintf("rules[input][%d][%s]", idx, "action"), rule.Action)
	}

	encodedParams := formParams.Encode()
	log.Println(encodedParams)

	_, err := c.makeAPICall("POST", fmt.Sprintf("%s/firewall/%s", c.url, firewall.IP), strings.NewReader(encodedParams), http.StatusAccepted)
	if err != nil {
		return err
	}

	return nil
}
