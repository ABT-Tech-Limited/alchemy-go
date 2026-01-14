// Package wallet provides high-level wallet operations using the Alchemy API.
package wallet

import (
	"github.com/ABT-Tech-Limited/alchemy-go/data"
	"github.com/ABT-Tech-Limited/alchemy-go/node"
)

// Client provides wallet-related operations.
type Client struct {
	data *data.Client
	node *node.Client
}

// NewClient creates a new Wallet client.
func NewClient(dataClient *data.Client, nodeClient *node.Client) *Client {
	return &Client{
		data: dataClient,
		node: nodeClient,
	}
}

// Data returns the underlying Data client.
func (c *Client) Data() *data.Client {
	return c.data
}

// Node returns the underlying Node client.
func (c *Client) Node() *node.Client {
	return c.node
}
