package hetzner

import (
	"strconv"
)

type StorageBoxClient struct {
	c *Client
}

func (c *Client) StorageBox() *StorageBoxClient {
	return &StorageBoxClient{
		c: c,
	}
}

func (s *StorageBoxClient) List() ([]*StorageBoxSummary, error) {
	var d []map[string]*StorageBoxSummary
	err := s.c.list("storagebox", &d)

	var res []*StorageBoxSummary
	for _, e := range d {
		res = append(res, e["storagebox"])
	}
	return res, err
}

func (s *StorageBoxClient) Info(id int) (*StorageBox, error) {
	var d map[string]*StorageBox
	if err := s.c.get("storagebox/"+strconv.Itoa(id), &d); err != nil {
		return nil, err
	}
	return d["storagebox"], nil
}

func (s *StorageBoxClient) Update(req *StorageBoxRequest) (*StorageBoxSummary, error) {
	var d map[string]*StorageBoxSummary
	err := s.c.post("storagebox/"+strconv.Itoa(req.ID), req, &d)
	if err != nil {
		return nil, err
	}
	return d["storagebox"], nil
}
