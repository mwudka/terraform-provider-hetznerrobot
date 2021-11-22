package hetznerrobot

// https://robot.your-server.de/doc/webservice/en.html#server

import (
	"encoding/json"
	"fmt"
	"go/types"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type TServerResponse struct {
	Server TServer `json:"server"`
}

type TServer struct {
	ServerIPv4       string      `json:"server_ip"`
	ServerIPv6       string      `json:"server_ipv6_net"`
	ServerID         int         `json:"server_number"`
	ServerName       string      `json:"server_name"`
	Product          string      `json:"product"`
	DataCenter       string      `json:"dc"`
	Traffic          string      `json:"traffic"`
	Status           string      `json:"status"`
	Cancelled        bool        `json:"cancelled"`
	PaidUntil        string      `json:"paid_until"`
	IPList           types.Array `json:"ip"`
	SubnetList       types.Array `json:"subnet"`
	ResetAvail       bool        `json:"reset"`
	RescueAvail      bool        `json:"rescue"`
	VNCAvail         bool        `json:"vnc"`
	WinAvail         bool        `json:"windows"`
	PleskAvail       bool        `json:"plesk"`
	CPanelAvail      bool        `json:"cpanel"`
	WOLAvail         bool        `json:"wol"`
	HotSwapAvail     bool        `json:"hot_swap"`
	LinkedStorageBox int16       `json:"storagebox"`
}

func (c *HetznerRobotClient) setServerName(serverID int, serverName string) (*TServer, error) {
	formParams := url.Values{}
	formParams.Set("server_name", serverName)
	encodedParams := formParams.Encode()
	log.Println(encodedParams)

	bytes, err := c.makeAPICall("POST", fmt.Sprintf("%s/server/%d", c.url, serverID), strings.NewReader(encodedParams), http.StatusAccepted)
	if err != nil {
		return nil, err
	}

	server := TServerResponse{}
	if err = json.Unmarshal(bytes, &server); err != nil {
		return nil, err
	}
	return &server.Server, nil
}

func (c *HetznerRobotClient) getServer(serverID int) (*TServer, error) {

	bytes, err := c.makeAPICall("GET", fmt.Sprintf("%s/server/%d", c.url, serverID), nil, http.StatusOK)
	if err != nil {
		return nil, err
	}

	server := TServerResponse{}
	if err = json.Unmarshal(bytes, &server); err != nil {
		return nil, err
	}
	return &server.Server, nil
}
