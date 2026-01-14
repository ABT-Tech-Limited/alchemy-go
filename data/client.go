// Package data provides the Data API client for Alchemy's enhanced APIs.
package data

import (
	"github.com/ABT-Tech-Limited/alchemy-go/client"
)

// Client is the Data API client.
type Client struct {
	http   *client.HTTPClient
	rpc    *client.JSONRPCClient
	nftURL string
}

// NewClient creates a new Data API client.
func NewClient(httpClient *client.HTTPClient, rpc *client.JSONRPCClient, nftURL string) *Client {
	return &Client{
		http:   httpClient,
		rpc:    rpc,
		nftURL: nftURL,
	}
}

// HTTP returns the underlying HTTP client.
func (c *Client) HTTP() *client.HTTPClient {
	return c.http
}

// RPC returns the underlying JSON-RPC client.
func (c *Client) RPC() *client.JSONRPCClient {
	return c.rpc
}
