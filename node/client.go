// Package node provides the Node API client for JSON-RPC methods.
package node

import (
	"github.com/ABT-Tech-Limited/alchemy-go/client"
)

// Client is the Node API client for making JSON-RPC calls.
type Client struct {
	rpc *client.JSONRPCClient
}

// NewClient creates a new Node API client.
func NewClient(rpc *client.JSONRPCClient) *Client {
	return &Client{
		rpc: rpc,
	}
}

// RPC returns the underlying JSON-RPC client.
func (c *Client) RPC() *client.JSONRPCClient {
	return c.rpc
}
