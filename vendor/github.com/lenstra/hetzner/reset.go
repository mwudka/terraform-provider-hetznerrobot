package hetzner

type ResetClient struct {
	c *Client
}

func (c *Client) Reset() *ResetClient {
	return &ResetClient{
		c: c,
	}
}

func (b *ResetClient) Reset(req *ResetRequest) (*Reset, error) {
	var d map[string]*Reset
	if err := b.c.post("reset/"+req.ServerIP, req, &d); err != nil {
		return nil, err
	}
	return d["reset"], nil
}
