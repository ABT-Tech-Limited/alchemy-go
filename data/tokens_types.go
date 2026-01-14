package data

import (
	"github.com/ABT-Tech-Limited/alchemy-go/types"
)

// TokenSpec represents the type of token to query.
type TokenSpec string

// Token specifications.
const (
	TokenSpecERC20       TokenSpec = "erc20"
	TokenSpecNativeToken TokenSpec = "NATIVE_TOKEN"
)

// TokenBalancesParams represents the parameters for getTokenBalances.
type TokenBalancesParams struct {
	// Address is the wallet address to query.
	Address types.Address `json:"address"`
	// TokenSpec is the type of tokens to query (optional).
	TokenSpec TokenSpec `json:"tokenSpec,omitempty"`
	// ContractAddresses is a list of specific contract addresses to query (optional).
	ContractAddresses []types.Address `json:"contractAddresses,omitempty"`
	// PageKey is the pagination key for fetching more results.
	PageKey string `json:"pageKey,omitempty"`
	// MaxCount is the maximum number of results per page.
	MaxCount int `json:"maxCount,omitempty"`
}

// NewTokenBalancesParams creates a new TokenBalancesParams.
func NewTokenBalancesParams(address types.Address) *TokenBalancesParams {
	return &TokenBalancesParams{
		Address: address,
	}
}

// SetTokenSpec sets the token specification.
func (p *TokenBalancesParams) SetTokenSpec(spec TokenSpec) *TokenBalancesParams {
	p.TokenSpec = spec
	return p
}

// SetContractAddresses sets specific contract addresses to query.
func (p *TokenBalancesParams) SetContractAddresses(addresses []types.Address) *TokenBalancesParams {
	p.ContractAddresses = addresses
	return p
}

// SetMaxCount sets the maximum number of results.
func (p *TokenBalancesParams) SetMaxCount(count int) *TokenBalancesParams {
	p.MaxCount = count
	return p
}

// TokenBalancesResponse represents the response from getTokenBalances.
type TokenBalancesResponse struct {
	// Address is the queried address.
	Address types.Address `json:"address"`
	// TokenBalances is the list of token balances.
	TokenBalances []TokenBalance `json:"tokenBalances"`
	// PageKey is the pagination key for fetching more results.
	PageKey string `json:"pageKey,omitempty"`
}

// HasMore returns true if there are more results available.
func (r *TokenBalancesResponse) HasMore() bool {
	return r.PageKey != ""
}

// TokenBalance represents a single token balance.
type TokenBalance struct {
	// ContractAddress is the token contract address.
	ContractAddress types.Address `json:"contractAddress"`
	// TokenBalance is the balance in the smallest unit (hex encoded).
	TokenBalance *string `json:"tokenBalance,omitempty"`
	// Error is the error message if the balance couldn't be fetched.
	Error *string `json:"error,omitempty"`
}

// HasError returns true if there was an error fetching the balance.
func (b *TokenBalance) HasError() bool {
	return b.Error != nil
}

// TokenMetadata represents token metadata.
type TokenMetadata struct {
	// Name is the token name.
	Name *string `json:"name,omitempty"`
	// Symbol is the token symbol.
	Symbol *string `json:"symbol,omitempty"`
	// Decimals is the number of decimals.
	Decimals *int `json:"decimals,omitempty"`
	// Logo is the token logo URL.
	Logo *string `json:"logo,omitempty"`
}

// TokensForOwnerResponse represents the response from getTokensForOwner.
type TokensForOwnerResponse struct {
	// Tokens is the list of tokens owned by the address.
	Tokens []OwnedToken `json:"tokens"`
	// PageKey is the pagination key for fetching more results.
	PageKey string `json:"pageKey,omitempty"`
}

// HasMore returns true if there are more results available.
func (r *TokensForOwnerResponse) HasMore() bool {
	return r.PageKey != ""
}

// OwnedToken represents a token owned by an address.
type OwnedToken struct {
	// ContractAddress is the token contract address.
	ContractAddress types.Address `json:"contractAddress"`
	// RawBalance is the raw balance (hex encoded).
	RawBalance string `json:"rawBalance"`
	// Balance is the formatted balance.
	Balance string `json:"balance"`
	// Name is the token name.
	Name *string `json:"name,omitempty"`
	// Symbol is the token symbol.
	Symbol *string `json:"symbol,omitempty"`
	// Decimals is the number of decimals.
	Decimals *int `json:"decimals,omitempty"`
	// Logo is the token logo URL.
	Logo *string `json:"logo,omitempty"`
	// Error is the error message if there was an issue.
	Error *string `json:"error,omitempty"`
}

// TokenAllowanceParams represents the parameters for getTokenAllowance.
type TokenAllowanceParams struct {
	// Contract is the token contract address.
	Contract types.Address `json:"contract"`
	// Owner is the owner address.
	Owner types.Address `json:"owner"`
	// Spender is the spender address.
	Spender types.Address `json:"spender"`
}

// TokenAllowanceResponse represents the response from getTokenAllowance.
type TokenAllowanceResponse struct {
	// Allowance is the allowance amount (hex encoded).
	Allowance string `json:"allowance"`
}
