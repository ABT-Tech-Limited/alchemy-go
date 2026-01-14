// Package alchemy provides a Golang SDK for the Alchemy API.
//
// The SDK covers Node API (JSON-RPC), Data API (Transfers, NFTs, Tokens),
// and Wallet API for blockchain data access.
//
// Example usage:
//
//	client, err := alchemy.New(alchemy.Config{
//	    APIKey:  os.Getenv("ALCHEMY_API_KEY"),
//	    Network: alchemy.EthMainnet,
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Get block number
//	blockNum, _ := client.Node.BlockNumber(ctx)
//
//	// Get asset transfers
//	transfers, _ := client.Data.GetAssetTransfers(ctx, &data.AssetTransfersParams{
//	    ToAddress: &address,
//	    Category:  []data.AssetTransferCategory{data.CategoryERC20},
//	})
package alchemy

import (
	"github.com/ABT-Tech-Limited/alchemy-go/client"
	"github.com/ABT-Tech-Limited/alchemy-go/data"
	"github.com/ABT-Tech-Limited/alchemy-go/node"
	"github.com/ABT-Tech-Limited/alchemy-go/wallet"
)

// Alchemy is the main client for the Alchemy API.
type Alchemy struct {
	config *Config

	// Node provides access to JSON-RPC methods (eth_*, debug_*, etc.).
	Node *node.Client

	// Data provides access to enhanced data APIs (transfers, tokens, NFTs).
	Data *data.Client

	// Wallet provides high-level wallet operations.
	Wallet *wallet.Client
}

// New creates a new Alchemy client with the given configuration.
func New(cfg Config) (*Alchemy, error) {
	// Apply defaults
	cfg = cfg.WithDefaults()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Create HTTP client
	httpClient := client.NewHTTPClient(client.HTTPClientConfig{
		BaseURL:       cfg.GetBaseURL(),
		APIKey:        cfg.APIKey,
		Timeout:       cfg.Timeout,
		MaxRetries:    cfg.MaxRetries,
		RetryDelay:    cfg.RetryDelay,
		RetryMaxDelay: cfg.RetryMaxDelay,
		HTTPClient:    cfg.HTTPClient,
		Debug:         cfg.Debug,
	})

	// Create JSON-RPC client
	rpcClient := client.NewJSONRPCClient(httpClient)

	// Create sub-clients
	nodeClient := node.NewClient(rpcClient)
	dataClient := data.NewClient(httpClient, rpcClient, cfg.Network.NFTURL())
	walletClient := wallet.NewClient(dataClient, nodeClient)

	return &Alchemy{
		config: &cfg,
		Node:   nodeClient,
		Data:   dataClient,
		Wallet: walletClient,
	}, nil
}

// WithNetwork creates a new Alchemy client for a different network.
// This returns a new client instance; the original client is not modified.
func (a *Alchemy) WithNetwork(network Network) (*Alchemy, error) {
	cfg := *a.config
	cfg.Network = network
	cfg.BaseURL = "" // Reset to use network default
	return New(cfg)
}

// Network returns the current network.
func (a *Alchemy) Network() Network {
	return a.config.Network
}

// Config returns a copy of the current configuration.
func (a *Alchemy) Config() Config {
	return *a.config
}
