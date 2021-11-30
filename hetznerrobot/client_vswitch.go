package hetznerrobot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// https://robot.your-server.de/doc/webservice/en.html#vswitch

type TVSwitchResponse struct {
	id           int                             `json:"id"`
	name         string                          `json:"string"`
	vlan         int                             `json:"vlan"`
	isCancelled  bool                            `json:"cancelled"`
	server       []TVSwitchServersResponse       `json:"server"`
	subnet       []TVSwitchSubnetsResponse       `json:"subnet"`
	cloudNetwork []TVSwitchCloudNetworksResponse `json:"cloud_network"`
}

type TVSwitchServersResponse struct {
	ServerID   int    `json:"server_number"`
	ServerIPv4 string `json:"server_ip"`
	ServerIPv6 string `json:"server_ipv6_net"`
	Status     string `json:"status"`
}

type TVSwitchSubnetsResponse struct {
	Subnet  string `json:"ip"`
	Netmask string `json:"mask"`
	Gateway string `json:"gateway"`
}

type TVSwitchCloudNetworksResponse struct {
	ID      int    `json:"id"`
	Subnet  string `json:"ip"`
	Netmask string `json:"mask"`
	Gateway string `json:"gateway"`
}

func (c *HetznerRobotClient) createVSwitch(name string, vlan int) (*TVSwitchResponse, error) {
	formParams := url.Values{}
	formParams.Set("name", name)
	formParams.Set("vlan", strconv.Itoa(vlan))
	encodedParams := formParams.Encode()
	log.Println(encodedParams)

	r, err := http.NewRequest("POST", fmt.Sprintf("%s/vswitch", c.url), strings.NewReader(encodedParams))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
	if response.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("Hetzner API response HTTP %d: %s", response.StatusCode, bytes)
	}

	vSwitch := TVSwitchResponse{}
	if err = json.Unmarshal(bytes, &vSwitch); err != nil {
		return nil, err
	}
	return &vSwitch, nil
}

func (c *HetznerRobotClient) updateVSwitch(vSwitchID int, name string, vlan int) (*TVSwitchResponse, error) {
	formParams := url.Values{}
	formParams.Set("name", name)
	formParams.Set("vlan", strconv.Itoa(vlan))
	encodedParams := formParams.Encode()
	log.Println(encodedParams)

	r, err := http.NewRequest("POST", fmt.Sprintf("%s/vswitch/%d", c.url, vSwitchID), strings.NewReader(encodedParams))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
	if response.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("Hetzner API response HTTP %d: %s", response.StatusCode, bytes)
	}

	vSwitch := TVSwitchResponse{}
	if err = json.Unmarshal(bytes, &vSwitch); err != nil {
		return nil, err
	}
	return &vSwitch, nil
}

func (c *HetznerRobotClient) getVSwitch(vSwitchID int) (*TVSwitchResponse, error) {
	r, err := http.NewRequest("GET", fmt.Sprintf("%s/vswitch/%d", c.url, vSwitchID), nil)
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
	vSwitchResp := TVSwitchResponse{}
	if err = json.Unmarshal(bytes, &vSwitchResp); err != nil {
		return nil, err
	}

	return &vSwitchResp, nil
}

func (c *HetznerRobotClient) deleteVSwitch(vSwitchID int) error {
	r, err := http.NewRequest("DELETE", fmt.Sprintf("%s/vswitch/%d", c.url, vSwitchID), nil)
	if err != nil {
		return err
	}
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
	log.Printf("Hetzner response status %d\n%s", response.StatusCode, bytes)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Hetzner API response HTTP %d: %s", response.StatusCode, bytes)
	}

	return nil
}

func (c *HetznerRobotClient) addServerToVSwitch(vSwitchID int, serverID int) error {
	formParams := url.Values{}
	formParams.Set("server", strconv.Itoa(serverID))
	encodedParams := formParams.Encode()
	log.Println(encodedParams)

	r, err := http.NewRequest("POST", fmt.Sprintf("%s/vswitch/%d/server", c.url, vSwitchID), strings.NewReader(encodedParams))
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

func (c *HetznerRobotClient) removeServerFromVSwitch(vSwitchID int, serverID int) error {
	formParams := url.Values{}
	formParams.Set("server", strconv.Itoa(serverID))
	encodedParams := formParams.Encode()
	log.Println(encodedParams)

	r, err := http.NewRequest("DELETE", fmt.Sprintf("%s/vswitch/%d/server", c.url, vSwitchID), strings.NewReader(encodedParams))
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
