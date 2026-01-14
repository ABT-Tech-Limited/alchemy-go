package wallet

import (
	"context"
	"math/big"

	"github.com/ABT-Tech-Limited/alchemy-go/data"
	"github.com/ABT-Tech-Limited/alchemy-go/internal/hex"
	"github.com/ABT-Tech-Limited/alchemy-go/node"
	"github.com/ABT-Tech-Limited/alchemy-go/types"
)

// Balance represents a balance with both raw and formatted values.
type Balance struct {
	// Address is the wallet address.
	Address types.Address
	// Raw is the balance in the smallest unit (wei).
	Raw *big.Int
	// Formatted is the balance formatted in ETH (or native token).
	Formatted string
}

// GetBalance retrieves the native token balance for an address.
func (c *Client) GetBalance(ctx context.Context, address types.Address) (*Balance, error) {
	raw, err := c.node.GetBalance(ctx, address, node.BlockLatest)
	if err != nil {
		return nil, err
	}

	return &Balance{
		Address:   address,
		Raw:       raw,
		Formatted: formatWei(raw),
	}, nil
}

// GetBalanceAtBlock retrieves the native token balance at a specific block.
func (c *Client) GetBalanceAtBlock(ctx context.Context, address types.Address, block node.BlockNumberOrTag) (*Balance, error) {
	raw, err := c.node.GetBalance(ctx, address, block)
	if err != nil {
		return nil, err
	}

	return &Balance{
		Address:   address,
		Raw:       raw,
		Formatted: formatWei(raw),
	}, nil
}

// TokenBalancesResult represents the result of a token balances query.
type TokenBalancesResult struct {
	// Address is the wallet address.
	Address types.Address
	// Balances is the list of token balances.
	Balances []TokenBalanceInfo
	// PageKey is the pagination key for more results.
	PageKey string
}

// TokenBalanceInfo represents a single token balance with metadata.
type TokenBalanceInfo struct {
	// ContractAddress is the token contract address.
	ContractAddress types.Address
	// Balance is the raw balance (smallest unit).
	Balance *big.Int
	// BalanceFormatted is the formatted balance.
	BalanceFormatted string
	// Metadata contains token metadata (if available).
	Metadata *data.TokenMetadata
	// Error is any error that occurred fetching this balance.
	Error string
}

// GetTokenBalances retrieves ERC20 token balances for an address.
func (c *Client) GetTokenBalances(ctx context.Context, address types.Address, contractAddresses []types.Address) (*TokenBalancesResult, error) {
	params := data.NewTokenBalancesParams(address)
	if len(contractAddresses) > 0 {
		params.SetContractAddresses(contractAddresses)
	} else {
		params.SetTokenSpec(data.TokenSpecERC20)
	}

	resp, err := c.data.GetTokenBalances(ctx, params)
	if err != nil {
		return nil, err
	}

	result := &TokenBalancesResult{
		Address:  address,
		PageKey:  resp.PageKey,
		Balances: make([]TokenBalanceInfo, 0, len(resp.TokenBalances)),
	}

	for _, tb := range resp.TokenBalances {
		info := TokenBalanceInfo{
			ContractAddress: tb.ContractAddress,
		}

		if tb.Error != nil {
			info.Error = *tb.Error
		} else if tb.TokenBalance != nil {
			balance, _ := hex.DecodeBigInt(*tb.TokenBalance)
			info.Balance = balance
		}

		result.Balances = append(result.Balances, info)
	}

	return result, nil
}

// GetTokenBalancesWithMetadata retrieves token balances with metadata.
func (c *Client) GetTokenBalancesWithMetadata(ctx context.Context, address types.Address, contractAddresses []types.Address) (*TokenBalancesResult, error) {
	result, err := c.GetTokenBalances(ctx, address, contractAddresses)
	if err != nil {
		return nil, err
	}

	// Fetch metadata for each token
	for i := range result.Balances {
		if result.Balances[i].Error != "" {
			continue
		}

		metadata, err := c.data.GetTokenMetadata(ctx, result.Balances[i].ContractAddress)
		if err != nil {
			continue // Ignore metadata errors
		}

		result.Balances[i].Metadata = metadata

		// Format the balance using decimals
		if result.Balances[i].Balance != nil && metadata.Decimals != nil {
			result.Balances[i].BalanceFormatted = formatTokenBalance(result.Balances[i].Balance, *metadata.Decimals)
		}
	}

	return result, nil
}

// GetAllTokenBalances retrieves all ERC20 token balances with pagination.
func (c *Client) GetAllTokenBalances(ctx context.Context, address types.Address) (*TokenBalancesResult, error) {
	var allBalances []TokenBalanceInfo
	pageKey := ""

	for {
		params := data.NewTokenBalancesParams(address).
			SetTokenSpec(data.TokenSpecERC20)
		if pageKey != "" {
			params.PageKey = pageKey
		}

		resp, err := c.data.GetTokenBalances(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, tb := range resp.TokenBalances {
			info := TokenBalanceInfo{
				ContractAddress: tb.ContractAddress,
			}

			if tb.Error != nil {
				info.Error = *tb.Error
			} else if tb.TokenBalance != nil {
				balance, _ := hex.DecodeBigInt(*tb.TokenBalance)
				info.Balance = balance
			}

			allBalances = append(allBalances, info)
		}

		if !resp.HasMore() {
			break
		}
		pageKey = resp.PageKey
	}

	return &TokenBalancesResult{
		Address:  address,
		Balances: allBalances,
	}, nil
}

// formatWei formats a wei value as ETH.
func formatWei(wei *big.Int) string {
	if wei == nil {
		return "0"
	}

	// 1 ETH = 10^18 wei
	eth := new(big.Float).SetInt(wei)
	divisor := new(big.Float).SetInt(big.NewInt(1e18))
	eth.Quo(eth, divisor)

	return eth.Text('f', 18)
}

// formatTokenBalance formats a token balance using its decimals.
func formatTokenBalance(balance *big.Int, decimals int) string {
	if balance == nil {
		return "0"
	}

	if decimals == 0 {
		return balance.String()
	}

	// Create divisor: 10^decimals
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

	// Convert to float for division
	balanceFloat := new(big.Float).SetInt(balance)
	divisorFloat := new(big.Float).SetInt(divisor)
	result := new(big.Float).Quo(balanceFloat, divisorFloat)

	return result.Text('f', decimals)
}
