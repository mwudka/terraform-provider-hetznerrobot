package hetzner

import (
	"fmt"
)

type BootClient struct {
	c *Client
}

func (c *Client) Boot() *BootClient {
	return &BootClient{
		c: c,
	}
}

func (b *BootClient) Rescue(req *RescueRequest) (*Rescue, error) {
	var d map[string]*Rescue
	if err := b.c.post(fmt.Sprintf("boot/%s/rescue", req.ServerIP), req, &d); err != nil {
		return nil, err
	}
	return d["rescue"], nil
}
