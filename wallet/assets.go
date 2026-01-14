package wallet

import (
	"context"

	"github.com/ABT-Tech-Limited/alchemy-go/data"
	"github.com/ABT-Tech-Limited/alchemy-go/types"
)

// NFTsResult represents the result of an NFT query.
type NFTsResult struct {
	// Address is the wallet address.
	Address types.Address
	// NFTs is the list of owned NFTs.
	NFTs []data.OwnedNFT
	// TotalCount is the total number of NFTs.
	TotalCount int
	// PageKey is the pagination key for more results.
	PageKey string
}

// NFTQueryOptions provides options for NFT queries.
type NFTQueryOptions struct {
	// ContractAddresses filters NFTs by contract addresses.
	ContractAddresses []types.Address
	// ExcludeSpam excludes spam NFTs.
	ExcludeSpam bool
	// ExcludeAirdrops excludes airdropped NFTs.
	ExcludeAirdrops bool
	// WithMetadata includes NFT metadata.
	WithMetadata bool
	// PageSize is the number of results per page.
	PageSize int
}

// DefaultNFTQueryOptions returns default query options.
func DefaultNFTQueryOptions() *NFTQueryOptions {
	return &NFTQueryOptions{
		ExcludeSpam:  true,
		WithMetadata: true,
		PageSize:     100,
	}
}

// GetNFTs retrieves NFTs owned by an address.
func (c *Client) GetNFTs(ctx context.Context, address types.Address, options *NFTQueryOptions) (*NFTsResult, error) {
	if options == nil {
		options = DefaultNFTQueryOptions()
	}

	params := data.NewNFTsForOwnerParams(address)

	if len(options.ContractAddresses) > 0 {
		params.SetContractAddresses(options.ContractAddresses)
	}

	if options.ExcludeSpam || options.ExcludeAirdrops {
		var filters []data.NFTFilter
		if options.ExcludeSpam {
			filters = append(filters, data.NFTFilterSpam)
		}
		if options.ExcludeAirdrops {
			filters = append(filters, data.NFTFilterAirdrops)
		}
		params.SetExcludeFilters(filters)
	}

	params.SetWithMetadata(options.WithMetadata)

	if options.PageSize > 0 {
		params.SetPageSize(options.PageSize)
	}

	resp, err := c.data.GetNFTsForOwner(ctx, params)
	if err != nil {
		return nil, err
	}

	return &NFTsResult{
		Address:    address,
		NFTs:       resp.OwnedNFTs,
		TotalCount: resp.TotalCount,
		PageKey:    resp.PageKey,
	}, nil
}

// GetAllNFTs retrieves all NFTs owned by an address (handles pagination).
func (c *Client) GetAllNFTs(ctx context.Context, address types.Address, options *NFTQueryOptions) (*NFTsResult, error) {
	if options == nil {
		options = DefaultNFTQueryOptions()
	}

	var allNFTs []data.OwnedNFT
	totalCount := 0

	params := data.NewNFTsForOwnerParams(address)

	if len(options.ContractAddresses) > 0 {
		params.SetContractAddresses(options.ContractAddresses)
	}

	if options.ExcludeSpam || options.ExcludeAirdrops {
		var filters []data.NFTFilter
		if options.ExcludeSpam {
			filters = append(filters, data.NFTFilterSpam)
		}
		if options.ExcludeAirdrops {
			filters = append(filters, data.NFTFilterAirdrops)
		}
		params.SetExcludeFilters(filters)
	}

	params.SetWithMetadata(options.WithMetadata)

	if options.PageSize > 0 {
		params.SetPageSize(options.PageSize)
	}

	for {
		resp, err := c.data.GetNFTsForOwner(ctx, params)
		if err != nil {
			return nil, err
		}

		allNFTs = append(allNFTs, resp.OwnedNFTs...)
		totalCount = resp.TotalCount

		if !resp.HasMore() {
			break
		}
		params.PageKey = resp.PageKey
	}

	return &NFTsResult{
		Address:    address,
		NFTs:       allNFTs,
		TotalCount: totalCount,
	}, nil
}

// GetERC721Assets retrieves ERC721 NFTs owned by an address.
func (c *Client) GetERC721Assets(ctx context.Context, address types.Address) ([]data.OwnedNFT, error) {
	result, err := c.GetAllNFTs(ctx, address, DefaultNFTQueryOptions())
	if err != nil {
		return nil, err
	}

	var erc721s []data.OwnedNFT
	for _, nft := range result.NFTs {
		if nft.TokenType == string(data.NFTTokenTypeERC721) || nft.Contract.TokenType == string(data.NFTTokenTypeERC721) {
			erc721s = append(erc721s, nft)
		}
	}

	return erc721s, nil
}

// GetERC1155Assets retrieves ERC1155 NFTs owned by an address.
func (c *Client) GetERC1155Assets(ctx context.Context, address types.Address) ([]data.OwnedNFT, error) {
	result, err := c.GetAllNFTs(ctx, address, DefaultNFTQueryOptions())
	if err != nil {
		return nil, err
	}

	var erc1155s []data.OwnedNFT
	for _, nft := range result.NFTs {
		if nft.TokenType == string(data.NFTTokenTypeERC1155) || nft.Contract.TokenType == string(data.NFTTokenTypeERC1155) {
			erc1155s = append(erc1155s, nft)
		}
	}

	return erc1155s, nil
}

// AssetSummary provides a summary of all assets owned by an address.
type AssetSummary struct {
	// Address is the wallet address.
	Address types.Address
	// NativeBalance is the native token balance.
	NativeBalance *Balance
	// TokenCount is the number of ERC20 tokens held.
	TokenCount int
	// NFTCount is the number of NFTs held.
	NFTCount int
	// ERC721Count is the number of ERC721 NFTs.
	ERC721Count int
	// ERC1155Count is the number of ERC1155 NFTs.
	ERC1155Count int
}

// GetAssetSummary retrieves a summary of all assets for an address.
func (c *Client) GetAssetSummary(ctx context.Context, address types.Address) (*AssetSummary, error) {
	summary := &AssetSummary{
		Address: address,
	}

	// Get native balance
	balance, err := c.GetBalance(ctx, address)
	if err != nil {
		return nil, err
	}
	summary.NativeBalance = balance

	// Get token count
	tokens, err := c.GetTokenBalances(ctx, address, nil)
	if err != nil {
		return nil, err
	}
	// Count tokens with non-zero balance
	for _, tb := range tokens.Balances {
		if tb.Balance != nil && tb.Balance.Sign() > 0 {
			summary.TokenCount++
		}
	}

	// Get NFT summary
	nftParams := data.NewNFTsForOwnerParams(address).
		SetWithMetadata(false).
		SetPageSize(1)

	nftResp, err := c.data.GetNFTsForOwner(ctx, nftParams)
	if err != nil {
		return nil, err
	}
	summary.NFTCount = nftResp.TotalCount

	// Get detailed NFT breakdown if there are NFTs
	if summary.NFTCount > 0 {
		allNFTs, err := c.GetAllNFTs(ctx, address, &NFTQueryOptions{
			WithMetadata: false,
			ExcludeSpam:  true,
		})
		if err == nil {
			for _, nft := range allNFTs.NFTs {
				switch nft.TokenType {
				case string(data.NFTTokenTypeERC721):
					summary.ERC721Count++
				case string(data.NFTTokenTypeERC1155):
					summary.ERC1155Count++
				}
			}
		}
	}

	return summary, nil
}
