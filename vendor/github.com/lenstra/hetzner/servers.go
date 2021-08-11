package hetzner

type ServersClient struct {
	c *Client
}

func (c *Client) Servers() *ServersClient {
	return &ServersClient{c: c}
}

func (s *ServersClient) List() ([]*ServerSummary, error) {
	var d []map[string]*ServerSummary
	err := s.c.list("server", &d)

	var res []*ServerSummary
	for _, e := range d {
		res = append(res, e["server"])
	}
	return res, err
}

func (s *ServersClient) Info(ip string) (*Server, error) {
	var d map[string]*Server
	if err := s.c.get("server/"+ip, &d); err != nil {
		return nil, err
	}
	return d["server"], nil
}

func (s *ServersClient) Update(req *ServerRequest) (*Server, error) {
	var d map[string]*Server
	err := s.c.post("server/"+req.ServerIP, req, &d)
	if err != nil {
		return nil, err
	}
	return d["server"], nil
}
