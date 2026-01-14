// Package main demonstrates basic usage of the Alchemy SDK.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ABT-Tech-Limited/alchemy-go"
	"github.com/ABT-Tech-Limited/alchemy-go/types"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("ALCHEMY_API_KEY")
	if apiKey == "" {
		log.Fatal("ALCHEMY_API_KEY environment variable is required")
	}

	// Create Alchemy client
	client, err := alchemy.New(alchemy.Config{
		APIKey:  apiKey,
		Network: alchemy.EthMainnet,
		Timeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to create Alchemy client: %v", err)
	}

	ctx := context.Background()

	// Example 1: Get current block number
	fmt.Println("=== Block Number ===")
	blockNum, err := client.Node.BlockNumber(ctx)
	if err != nil {
		log.Printf("Failed to get block number: %v", err)
	} else {
		fmt.Printf("Current block number: %d\n", blockNum)
	}

	// Example 2: Get chain ID
	fmt.Println("\n=== Chain ID ===")
	chainID, err := client.Node.ChainID(ctx)
	if err != nil {
		log.Printf("Failed to get chain ID: %v", err)
	} else {
		fmt.Printf("Chain ID: %d\n", chainID)
	}

	// Example 3: Get gas price
	fmt.Println("\n=== Gas Price ===")
	gasPrice, err := client.Node.GasPrice(ctx)
	if err != nil {
		log.Printf("Failed to get gas price: %v", err)
	} else {
		// Convert wei to gwei
		gwei := gasPrice.Int64() / 1e9
		fmt.Printf("Gas price: %d gwei\n", gwei)
	}

	// Example 4: Get latest block
	fmt.Println("\n=== Latest Block ===")
	block, err := client.Node.GetBlockByNumber(ctx, "latest", false)
	if err != nil {
		log.Printf("Failed to get block: %v", err)
	} else {
		fmt.Printf("Block hash: %s\n", block.Hash)
		fmt.Printf("Block number: %d\n", block.Number.Uint64())
		fmt.Printf("Timestamp: %d\n", block.Timestamp.Uint64())
		fmt.Printf("Transaction count: %d\n", block.TransactionCount())
	}

	// Example 5: Get balance for Vitalik's address
	fmt.Println("\n=== Balance Check ===")
	vitalikAddr := types.MustParseAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	balance, err := client.Wallet.GetBalance(ctx, vitalikAddr)
	if err != nil {
		log.Printf("Failed to get balance: %v", err)
	} else {
		fmt.Printf("Address: %s\n", balance.Address)
		fmt.Printf("Balance: %s ETH\n", balance.Formatted[:20]) // Truncate for display
	}

	fmt.Println("\n=== Done ===")
}
