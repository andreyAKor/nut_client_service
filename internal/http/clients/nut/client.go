package nut

import (
	"github.com/pkg/errors"
	nut "github.com/robbiet480/go.nut"
)

type Client struct {
	host     string
	port     int
	username string
	password string
}

func New(host string, port int, username, password string) (*Client, error) {
	return &Client{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}, nil
}

func (c *Client) GetUPSList() ([]nut.UPS, error) {
	client, err := c.connect()
	if err != nil {
		return nil, errors.Wrap(err, "connect fail")
	}

	list, err := client.GetUPSList()
	if err != nil {
		return nil, errors.Wrap(err, "get UPS list fail")
	}

	if err := c.disconnect(client); err != nil {
		return nil, errors.Wrap(err, "disconnect fail")
	}

	return list, nil
}

func (c *Client) connect() (*nut.Client, error) {
	client, err := nut.Connect(c.host, c.port)
	if err != nil {
		return nil, errors.Wrap(err, "connect fail")
	}

	if len(c.username) > 0 || len(c.password) > 0 {
		ok, err := client.Authenticate(c.username, c.password)
		if err != nil {
			return nil, errors.Wrap(err, "authenticate fail")
		}
		if !ok {
			return nil, errors.Wrap(err, "authenticate error")
		}
	}

	return &client, nil
}

func (c *Client) disconnect(client *nut.Client) error {
	ok, err := client.Disconnect()
	if err != nil {
		return errors.Wrap(err, "disconnect fail")
	}
	if !ok {
		return errors.Wrap(err, "disconnect error")
	}

	return nil
}
