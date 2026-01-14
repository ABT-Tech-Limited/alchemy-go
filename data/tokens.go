package data

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ABT-Tech-Limited/alchemy-go/types"
)

// GetTokenBalances retrieves token balances for an address.
func (c *Client) GetTokenBalances(ctx context.Context, params *TokenBalancesParams) (*TokenBalancesResponse, error) {
	// Build the request parameters
	reqParams := make([]interface{}, 0, 3)
	reqParams = append(reqParams, params.Address.String())

	// Add contract addresses or token spec
	if len(params.ContractAddresses) > 0 {
		addrs := make([]string, len(params.ContractAddresses))
		for i, addr := range params.ContractAddresses {
			addrs[i] = addr.String()
		}
		reqParams = append(reqParams, addrs)
	} else if params.TokenSpec != "" {
		reqParams = append(reqParams, string(params.TokenSpec))
	}

	// Add options if needed
	if params.PageKey != "" || params.MaxCount > 0 {
		options := make(map[string]interface{})
		if params.PageKey != "" {
			options["pageKey"] = params.PageKey
		}
		if params.MaxCount > 0 {
			options["maxCount"] = params.MaxCount
		}
		reqParams = append(reqParams, options)
	}

	var result TokenBalancesResponse
	if err := c.rpc.Call(ctx, "alchemy_getTokenBalances", reqParams, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTokenBalancesForAddresses retrieves token balances for multiple addresses.
// This is a convenience method that makes multiple calls.
func (c *Client) GetTokenBalancesForAddresses(ctx context.Context, addresses []types.Address, contractAddresses []types.Address) (map[types.Address]*TokenBalancesResponse, error) {
	results := make(map[types.Address]*TokenBalancesResponse)

	for _, addr := range addresses {
		params := NewTokenBalancesParams(addr)
		if len(contractAddresses) > 0 {
			params.SetContractAddresses(contractAddresses)
		}

		result, err := c.GetTokenBalances(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to get token balances for %s: %w", addr, err)
		}
		results[addr] = result
	}

	return results, nil
}

// GetTokenMetadata retrieves metadata for a token.
func (c *Client) GetTokenMetadata(ctx context.Context, contractAddress types.Address) (*TokenMetadata, error) {
	var result TokenMetadata
	if err := c.rpc.Call(ctx, "alchemy_getTokenMetadata", []interface{}{contractAddress.String()}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTokensForOwner retrieves all tokens owned by an address.
func (c *Client) GetTokensForOwner(ctx context.Context, owner types.Address, pageKey string) (*TokensForOwnerResponse, error) {
	params := make(map[string]interface{})
	params["owner"] = owner.String()
	if pageKey != "" {
		params["pageKey"] = pageKey
	}

	var result TokensForOwnerResponse
	if err := c.rpc.Call(ctx, "alchemy_getTokensForOwner", []interface{}{params}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTokenAllowance retrieves the allowance for a spender.
func (c *Client) GetTokenAllowance(ctx context.Context, params *TokenAllowanceParams) (*TokenAllowanceResponse, error) {
	reqParams := map[string]string{
		"contract": params.Contract.String(),
		"owner":    params.Owner.String(),
		"spender":  params.Spender.String(),
	}

	var result TokenAllowanceResponse
	if err := c.rpc.Call(ctx, "alchemy_getTokenAllowance", []interface{}{reqParams}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// nftHTTPGet makes a GET request to the NFT API.
func (c *Client) nftHTTPGet(ctx context.Context, path string, query url.Values, result interface{}) error {
	// Build the URL
	fullURL := c.nftURL + "/" + path
	if len(query) > 0 {
		fullURL = fullURL + "?" + query.Encode()
	}

	body, err := c.http.Get(ctx, fullURL)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, result)
}
