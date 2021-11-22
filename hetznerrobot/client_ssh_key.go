package hetznerrobot

// https://robot.your-server.de/doc/webservice/en.html#key

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type SSHKeyResponse struct {
	SSHKey SSHKey `json:"key"`
}

type SSHKey struct {
	Name        string `json:"name"`
	FingerPrint string `json:"fingerprint"`
	Type        string `json:"type"`
	Size        int    `json:"size"`
	Data        string `json:"data"`
}

func (c *HetznerRobotClient) getSSHKey(fingerprint string) (*SSHKey, error) {
	r, err := http.NewRequest("GET", fmt.Sprintf("%s/key/%s", c.url, fingerprint), nil)
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
	sshKey := SSHKeyResponse{}
	if err = json.Unmarshal(bytes, &sshKey); err != nil {
		return nil, err
	}
	return &sshKey.SSHKey, nil
}

func (c *HetznerRobotClient) addSSHKey(name string, data string) (*SSHKey, error) {
	formParams := url.Values{}
	formParams.Set("name", name)
	formParams.Set("data", data)
	encodedParams := formParams.Encode()
	log.Println(encodedParams)

	r, err := http.NewRequest("POST", fmt.Sprintf("%s/key", c.url), strings.NewReader(encodedParams))

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

	sshKey := SSHKeyResponse{}
	if err = json.Unmarshal(bytes, &sshKey); err != nil {
		return nil, err
	}
	return &sshKey.SSHKey, nil
}

func (c *HetznerRobotClient) updSSHKey(name string, fingerprint string) error {
	formParams := url.Values{}
	formParams.Set("name", name)
	encodedParams := formParams.Encode()
	log.Println(encodedParams)

	r, err := http.NewRequest("POST", fmt.Sprintf("%s/key/%s", c.url, fingerprint), strings.NewReader(encodedParams))

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
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Hetzner API response HTTP %d: %s", response.StatusCode, bytes)
	}

	return nil
}

func (c *HetznerRobotClient) delSSHKey(fingerprint string) error {
	r, err := http.NewRequest("DELETE", fmt.Sprintf("%s/key/%s", c.url, fingerprint), nil)
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
