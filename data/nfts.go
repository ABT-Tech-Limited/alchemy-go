package data

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"github.com/ABT-Tech-Limited/alchemy-go/types"
)

// GetNFTsForOwner retrieves NFTs owned by an address.
func (c *Client) GetNFTsForOwner(ctx context.Context, params *NFTsForOwnerParams) (*NFTsForOwnerResponse, error) {
	query := url.Values{}
	query.Set("owner", params.Owner.String())

	if len(params.ContractAddresses) > 0 {
		for _, addr := range params.ContractAddresses {
			query.Add("contractAddresses[]", addr.String())
		}
	}

	if params.WithMetadata != nil {
		query.Set("withMetadata", fmt.Sprintf("%t", *params.WithMetadata))
	}

	if params.OrderBy != "" {
		query.Set("orderBy", string(params.OrderBy))
	}

	for _, filter := range params.ExcludeFilters {
		query.Add("excludeFilters[]", string(filter))
	}

	for _, filter := range params.IncludeFilters {
		query.Add("includeFilters[]", string(filter))
	}

	if params.SpamConfidenceLevel != "" {
		query.Set("spamConfidenceLevel", string(params.SpamConfidenceLevel))
	}

	if params.TokenURITimeoutInMs != nil {
		query.Set("tokenUriTimeoutInMs", fmt.Sprintf("%d", *params.TokenURITimeoutInMs))
	}

	if params.PageKey != "" {
		query.Set("pageKey", params.PageKey)
	}

	if params.PageSize != nil {
		query.Set("pageSize", fmt.Sprintf("%d", *params.PageSize))
	}

	var result NFTsForOwnerResponse
	if err := c.nftGet(ctx, "getNFTsForOwner", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetNFTsForOwnerIterator returns an iterator for paginating through NFTs.
func (c *Client) GetNFTsForOwnerIterator(ctx context.Context, params *NFTsForOwnerParams) *NFTsForOwnerIterator {
	paramsCopy := *params
	return &NFTsForOwnerIterator{
		client: c,
		params: &paramsCopy,
		ctx:    ctx,
	}
}

// NFTsForOwnerIterator iterates through NFTs with pagination.
type NFTsForOwnerIterator struct {
	client  *Client
	params  *NFTsForOwnerParams
	ctx     context.Context
	current *NFTsForOwnerResponse
	index   int
	done    bool
	err     error
	mu      sync.Mutex
}

// Next returns the next NFT in the iteration.
func (it *NFTsForOwnerIterator) Next() (*OwnedNFT, error) {
	it.mu.Lock()
	defer it.mu.Unlock()

	if it.err != nil {
		return nil, it.err
	}

	if it.done {
		return nil, nil
	}

	if it.current == nil {
		if err := it.fetchNext(); err != nil {
			it.err = err
			return nil, err
		}
	}

	if it.index < len(it.current.OwnedNFTs) {
		nft := &it.current.OwnedNFTs[it.index]
		it.index++
		return nft, nil
	}

	if !it.current.HasMore() {
		it.done = true
		return nil, nil
	}

	it.params.PageKey = it.current.PageKey
	if err := it.fetchNext(); err != nil {
		it.err = err
		return nil, err
	}

	if len(it.current.OwnedNFTs) == 0 {
		it.done = true
		return nil, nil
	}

	nft := &it.current.OwnedNFTs[0]
	it.index = 1
	return nft, nil
}

// HasNext returns true if there are more NFTs to iterate.
func (it *NFTsForOwnerIterator) HasNext() bool {
	it.mu.Lock()
	defer it.mu.Unlock()

	if it.done || it.err != nil {
		return false
	}

	if it.current != nil && it.index < len(it.current.OwnedNFTs) {
		return true
	}

	if it.current != nil {
		return it.current.HasMore()
	}

	return true
}

// Error returns any error encountered during iteration.
func (it *NFTsForOwnerIterator) Error() error {
	it.mu.Lock()
	defer it.mu.Unlock()
	return it.err
}

// TotalCount returns the total count of NFTs (available after first fetch).
func (it *NFTsForOwnerIterator) TotalCount() int {
	it.mu.Lock()
	defer it.mu.Unlock()
	if it.current != nil {
		return it.current.TotalCount
	}
	return 0
}

// Collect returns all remaining NFTs as a slice.
func (it *NFTsForOwnerIterator) Collect() ([]OwnedNFT, error) {
	var nfts []OwnedNFT

	for {
		nft, err := it.Next()
		if err != nil {
			return nil, err
		}
		if nft == nil {
			break
		}
		nfts = append(nfts, *nft)
	}

	return nfts, nil
}

func (it *NFTsForOwnerIterator) fetchNext() error {
	result, err := it.client.GetNFTsForOwner(it.ctx, it.params)
	if err != nil {
		return err
	}
	it.current = result
	it.index = 0
	return nil
}

// GetNFTMetadata retrieves metadata for a specific NFT.
func (c *Client) GetNFTMetadata(ctx context.Context, params *NFTMetadataParams) (*OwnedNFT, error) {
	query := url.Values{}
	query.Set("contractAddress", params.ContractAddress.String())
	query.Set("tokenId", params.TokenID)

	if params.TokenType != nil {
		query.Set("tokenType", *params.TokenType)
	}

	if params.RefreshCache {
		query.Set("refreshCache", "true")
	}

	var result OwnedNFT
	if err := c.nftGet(ctx, "getNFTMetadata", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetContractMetadata retrieves metadata for an NFT contract.
func (c *Client) GetContractMetadata(ctx context.Context, contractAddress types.Address) (*NFTContractMetadata, error) {
	query := url.Values{}
	query.Set("contractAddress", contractAddress.String())

	var result NFTContractMetadata
	if err := c.nftGet(ctx, "getContractMetadata", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetNFTsForContract retrieves all NFTs for a contract.
func (c *Client) GetNFTsForContract(ctx context.Context, contractAddress types.Address, pageKey string, withMetadata bool) (*NFTsForContractResponse, error) {
	query := url.Values{}
	query.Set("contractAddress", contractAddress.String())
	query.Set("withMetadata", fmt.Sprintf("%t", withMetadata))

	if pageKey != "" {
		query.Set("startToken", pageKey)
	}

	var result NFTsForContractResponse
	if err := c.nftGet(ctx, "getNFTsForContract", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// NFTsForContractResponse represents the response from getNFTsForContract.
type NFTsForContractResponse struct {
	// NFTs is the list of NFTs in the contract.
	NFTs []OwnedNFT `json:"nfts"`
	// PageKey is the pagination key.
	PageKey string `json:"pageKey,omitempty"`
}

// HasMore returns true if there are more results available.
func (r *NFTsForContractResponse) HasMore() bool {
	return r.PageKey != ""
}

// GetOwnersForNFT retrieves owners for a specific NFT.
func (c *Client) GetOwnersForNFT(ctx context.Context, contractAddress types.Address, tokenID string) (*OwnersForNFTResponse, error) {
	query := url.Values{}
	query.Set("contractAddress", contractAddress.String())
	query.Set("tokenId", tokenID)

	var result OwnersForNFTResponse
	if err := c.nftGet(ctx, "getOwnersForNFT", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// OwnersForNFTResponse represents the response from getOwnersForNFT.
type OwnersForNFTResponse struct {
	// Owners is the list of owner addresses.
	Owners []types.Address `json:"owners"`
	// PageKey is the pagination key.
	PageKey string `json:"pageKey,omitempty"`
}

// GetOwnersForContract retrieves all owners for an NFT contract.
func (c *Client) GetOwnersForContract(ctx context.Context, contractAddress types.Address, pageKey string, withTokenBalances bool) (*OwnersForContractResponse, error) {
	query := url.Values{}
	query.Set("contractAddress", contractAddress.String())
	query.Set("withTokenBalances", fmt.Sprintf("%t", withTokenBalances))

	if pageKey != "" {
		query.Set("pageKey", pageKey)
	}

	var result OwnersForContractResponse
	if err := c.nftGet(ctx, "getOwnersForContract", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// OwnersForContractResponse represents the response from getOwnersForContract.
type OwnersForContractResponse struct {
	// Owners is the list of owners with their balances.
	Owners []ContractOwner `json:"ownerAddresses"`
	// PageKey is the pagination key.
	PageKey string `json:"pageKey,omitempty"`
}

// ContractOwner represents an owner of NFTs in a contract.
type ContractOwner struct {
	// OwnerAddress is the owner's address.
	OwnerAddress types.Address `json:"ownerAddress"`
	// TokenBalances contains the tokens owned (if withTokenBalances is true).
	TokenBalances []TokenBalanceEntry `json:"tokenBalances,omitempty"`
}

// TokenBalanceEntry represents a token balance entry.
type TokenBalanceEntry struct {
	// TokenID is the token ID.
	TokenID string `json:"tokenId"`
	// Balance is the balance.
	Balance string `json:"balance"`
}

// IsSpamContract checks if a contract is classified as spam.
func (c *Client) IsSpamContract(ctx context.Context, contractAddress types.Address) (bool, error) {
	query := url.Values{}
	query.Set("contractAddress", contractAddress.String())

	var result struct {
		IsSpamContract bool `json:"isSpamContract"`
	}
	if err := c.nftGet(ctx, "isSpamContract", query, &result); err != nil {
		return false, err
	}
	return result.IsSpamContract, nil
}

// nftGet makes a GET request to the NFT API endpoint.
func (c *Client) nftGet(ctx context.Context, method string, query url.Values, result interface{}) error {
	// Build the full URL: nftURL/apiKey/method
	baseURL := c.nftURL
	apiKey := c.http.BaseURL()

	// Extract API key from base URL (the part after the last /)
	if len(apiKey) > 0 {
		for i := len(apiKey) - 1; i >= 0; i-- {
			if apiKey[i] == '/' {
				apiKey = apiKey[i+1:]
				break
			}
		}
	}

	fullURL := baseURL + "/" + apiKey + "/" + method
	if len(query) > 0 {
		fullURL = fullURL + "?" + query.Encode()
	}

	body, err := c.http.Get(ctx, fullURL)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, result)
}
