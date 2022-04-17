package nut

import (
	"context"
	"time"

	nut_client "github.com/andreyAKor/nut_client"
	"github.com/pkg/errors"
)

const timeout = 1

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

// GetUPSList Returns a list of all UPSes provided by this NUT instance.
func (c *Client) GetUPSList(ctx context.Context) ([]*nut_client.UPS, error) {
	client, err := c.connect(ctx)
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

// SendCommand Sends a command to the UPS.
func (c *Client) SendCommand(ctx context.Context, name, command string) error {
	client, err := c.connect(ctx)
	if err != nil {
		return errors.Wrap(err, "connect fail")
	}

	list, err := client.GetUPSList()
	if err != nil {
		return errors.Wrap(err, "get UPS list fail")
	}
	for _, ups := range list {
		if ups.Name == name {
			ok, err := ups.SendCommand(command)
			if err != nil {
				return errors.Wrapf(err, `send command "%s" to UPS "%s" has failed`, name, command)
			}
			if !ok {
				return errors.Wrapf(err, `send command "%s" to UPS "%s" is error`, name, command)
			}
		}
	}

	if err := c.disconnect(client); err != nil {
		return errors.Wrap(err, "disconnect fail")
	}

	return nil
}

// connect Connecting to NUT.
func (c *Client) connect(ctx context.Context) (*nut_client.Client, error) {
	ctx, _ = context.WithTimeout(ctx, time.Second*timeout)

	client, err := nut_client.NewClient(ctx, c.host, c.port)
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

	return client, nil
}

// disconnect Gracefully disconnects from NUT by sending the LOGOUT command.
func (c *Client) disconnect(client *nut_client.Client) error {
	ok, err := client.Disconnect()
	if err != nil {
		return errors.Wrap(err, "disconnect fail")
	}
	if !ok {
		return errors.Wrap(err, "disconnect error")
	}

	return nil
}
