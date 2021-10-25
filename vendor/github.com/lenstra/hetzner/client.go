package hetzner

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/form"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-hclog"
)

func Bool(b bool) *bool {
	return &b
}

func Int(i int) *int {
	return &i
}

func String(s string) *string {
	return &s
}

type Config struct {
	Address  string
	Logger   hclog.Logger
	User     string
	Password string
}

type Client struct {
	address, user, password string
	client                  *http.Client
	logger                  hclog.Logger
}

func DefaultConfig() *Config {
	return &Config{
		Address: "https://robot-ws.your-server.de/",
	}
}

func NewClient(config *Config) *Client {
	c := &Client{
		client: cleanhttp.DefaultClient(),
	}

	if config == nil {
		config = DefaultConfig()
	}

	c.address = config.Address
	c.user = config.User
	c.password = config.Password

	if config.Logger == nil {
		c.logger = hclog.NewNullLogger()
	} else {
		c.logger = config.Logger
	}

	if !strings.HasSuffix(c.address, "/") {
		c.address += "/"
	}

	return c
}

func (c *Client) do(method, url string, body url.Values) ([]byte, error) {
	url = c.address + url

	var reader io.Reader
	if body != nil {
		reader = strings.NewReader(body.Encode())
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.SetBasicAuth(c.user, c.password)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("Got response", "method", method, "url", url, "status", resp.Status, "content", string(content))

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%s: %s", resp.Status, content)
	}
	return content, nil
}

func (c *Client) list(url string, d interface{}) error {
	content, err := c.do("GET", url, nil)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, d)
}

func (c *Client) get(url string, d interface{}) error {
	content, err := c.do("GET", url, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, d)
}

func (c *Client) post(u string, req interface{}, d interface{}) error {
	encoder := form.NewEncoder()
	body, err := encoder.Encode(req)
	if err != nil {
		return err
	}

	content, err := c.do("POST", u, body)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, d)
}

func (c *Client) delete(url string, d interface{}) error {
	content, err := c.do("DELETE", url, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, d)
}
