package hetznerrobot

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type HetznerRobotClient struct {
	username string
	password string
	url      string
}

func NewHetznerRobotClient(username string, password string, url string) HetznerRobotClient {
	return HetznerRobotClient{
		username: username,
		password: password,
		url:      url,
	}
}

func (c *HetznerRobotClient) makeAPICall(method string, uri string, body io.Reader, expectedStatusCode int) ([]byte, error) {
	r, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
	if response.StatusCode != expectedStatusCode {
		return nil, fmt.Errorf("Hetzner API response HTTP %d: %s", response.StatusCode, bytes)
	}

	return bytes, nil
}
