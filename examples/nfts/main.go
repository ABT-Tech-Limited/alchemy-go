// Package main demonstrates querying NFTs using the Alchemy SDK.
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

	// Query NFTs for owner
	fmt.Println("=== NFTs for Owner ===")
	fmt.Printf("Address: %s\n\n", address)

	params := data.NewNFTsForOwnerParams(address).
		SetWithMetadata(true).
		SetExcludeFilters([]data.NFTFilter{data.NFTFilterSpam}).
		SetPageSize(10)

	resp, err := client.Data.GetNFTsForOwner(ctx, params)
	if err != nil {
		log.Fatalf("Failed to get NFTs: %v", err)
	}

	fmt.Printf("Total NFTs: %d\n", resp.TotalCount)
	fmt.Printf("Returned: %d\n\n", len(resp.OwnedNFTs))

	for i, nft := range resp.OwnedNFTs {
		fmt.Printf("NFT #%d:\n", i+1)
		fmt.Printf("  Contract: %s\n", nft.Contract.Address)
		fmt.Printf("  Token ID: %s\n", nft.TokenID)
		fmt.Printf("  Type: %s\n", nft.TokenType)

		if nft.Contract.Name != nil {
			fmt.Printf("  Collection: %s\n", *nft.Contract.Name)
		}

		if nft.Name != nil {
			fmt.Printf("  Name: %s\n", *nft.Name)
		}

		if nft.Image != nil && nft.Image.ThumbnailURL != nil {
			fmt.Printf("  Thumbnail: %s\n", *nft.Image.ThumbnailURL)
		}

		fmt.Println()
	}

	if resp.HasMore() {
		fmt.Printf("More results available (pageKey: %s)\n", resp.PageKey)
	}

	// Example: Using wallet API for asset summary
	fmt.Println("\n=== Asset Summary ===")

	summary, err := client.Wallet.GetAssetSummary(ctx, address)
	if err != nil {
		log.Printf("Failed to get asset summary: %v", err)
	} else {
		fmt.Printf("Address: %s\n", summary.Address)
		fmt.Printf("Native Balance: %s ETH\n", summary.NativeBalance.Formatted[:20])
		fmt.Printf("ERC20 Tokens: %d\n", summary.TokenCount)
		fmt.Printf("Total NFTs: %d\n", summary.NFTCount)
		fmt.Printf("  ERC721: %d\n", summary.ERC721Count)
		fmt.Printf("  ERC1155: %d\n", summary.ERC1155Count)
	}

	// Example: Get specific NFT metadata
	fmt.Println("\n=== NFT Metadata ===")

	// BAYC #1
	baycContract := types.MustParseAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D")
	nftParams := data.NewNFTMetadataParams(baycContract, "1")

	nftMetadata, err := client.Data.GetNFTMetadata(ctx, nftParams)
	if err != nil {
		log.Printf("Failed to get NFT metadata: %v", err)
	} else {
		fmt.Printf("Contract: %s\n", nftMetadata.Contract.Address)
		fmt.Printf("Token ID: %s\n", nftMetadata.TokenID)
		if nftMetadata.Name != nil {
			fmt.Printf("Name: %s\n", *nftMetadata.Name)
		}
		if nftMetadata.Contract.Name != nil {
			fmt.Printf("Collection: %s\n", *nftMetadata.Contract.Name)
		}
	}

	fmt.Println("\n=== Done ===")
}
