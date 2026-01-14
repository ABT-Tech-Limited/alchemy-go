// Package main demonstrates querying asset transfers using the Alchemy SDK.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ABT-Tech-Limited/alchemy-go"
	"github.com/ABT-Tech-Limited/alchemy-go/data"
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

	// Example address (Vitalik's address)
	address := types.MustParseAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")

	// Query ERC20 incoming transfers
	fmt.Println("=== ERC20 Incoming Transfers ===")
	fmt.Printf("Address: %s\n\n", address)

	params := &data.AssetTransfersParams{
		ToAddress:        &address,
		Category:         []data.AssetTransferCategory{data.CategoryERC20},
		WithMetadata:     true,
		ExcludeZeroValue: true,
		Order:            data.SortDesc,
		MaxCount:         "0xa", // 10 results
	}

	resp, err := client.Data.GetAssetTransfers(ctx, params)
	if err != nil {
		log.Fatalf("Failed to get asset transfers: %v", err)
	}

	fmt.Printf("Found %d transfers\n\n", len(resp.Transfers))

	for i, transfer := range resp.Transfers {
		fmt.Printf("Transfer #%d:\n", i+1)
		fmt.Printf("  Hash: %s\n", transfer.Hash)
		fmt.Printf("  From: %s\n", transfer.From)
		fmt.Printf("  Block: %s\n", transfer.BlockNum)

		if transfer.Asset != nil {
			fmt.Printf("  Asset: %s\n", *transfer.Asset)
		}

		if transfer.Value != nil {
			fmt.Printf("  Value: %.6f\n", *transfer.Value)
		}

		if transfer.Metadata != nil {
			fmt.Printf("  Timestamp: %s\n", transfer.Metadata.BlockTimestamp)
		}

		fmt.Println()
	}

	if resp.HasMore() {
		fmt.Printf("More results available (pageKey: %s)\n", resp.PageKey)
	}

	// Example: Using iterator for pagination
	fmt.Println("\n=== Using Iterator ===")

	iterParams := &data.AssetTransfersParams{
		ToAddress:        &address,
		Category:         []data.AssetTransferCategory{data.CategoryExternal},
		WithMetadata:     true,
		ExcludeZeroValue: true,
		Order:            data.SortDesc,
		MaxCount:         "0x5", // 5 per page
	}

	iterator := client.Data.GetAssetTransfersIterator(ctx, iterParams)

	// Collect up to 10 transfers
	transfers, err := iterator.CollectN(10)
	if err != nil {
		log.Printf("Failed to collect transfers: %v", err)
	} else {
		fmt.Printf("Collected %d ETH transfers using iterator\n", len(transfers))
		for _, t := range transfers {
			if t.Value != nil {
				fmt.Printf("  %s: %.4f ETH\n", t.Hash[:16], *t.Value)
			}
		}
	}

	fmt.Println("\n=== Done ===")
}
