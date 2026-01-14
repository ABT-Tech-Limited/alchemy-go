package node

import (
	"context"
	"math/big"

	"github.com/ABT-Tech-Limited/alchemy-go/internal/hex"
	"github.com/ABT-Tech-Limited/alchemy-go/types"
)

// BlockNumber returns the current block number.
func (c *Client) BlockNumber(ctx context.Context) (uint64, error) {
	var result types.Quantity
	if err := c.rpc.Call(ctx, "eth_blockNumber", nil, &result); err != nil {
		return 0, err
	}
	return result.Uint64(), nil
}

// ChainID returns the chain ID.
func (c *Client) ChainID(ctx context.Context) (uint64, error) {
	var result types.Quantity
	if err := c.rpc.Call(ctx, "eth_chainId", nil, &result); err != nil {
		return 0, err
	}
	return result.Uint64(), nil
}

// GasPrice returns the current gas price in wei.
func (c *Client) GasPrice(ctx context.Context) (*big.Int, error) {
	var result types.Quantity
	if err := c.rpc.Call(ctx, "eth_gasPrice", nil, &result); err != nil {
		return nil, err
	}
	return result.BigInt(), nil
}

// MaxPriorityFeePerGas returns the current max priority fee per gas.
func (c *Client) MaxPriorityFeePerGas(ctx context.Context) (*big.Int, error) {
	var result types.Quantity
	if err := c.rpc.Call(ctx, "eth_maxPriorityFeePerGas", nil, &result); err != nil {
		return nil, err
	}
	return result.BigInt(), nil
}

// BlobBaseFee returns the current blob base fee.
func (c *Client) BlobBaseFee(ctx context.Context) (*big.Int, error) {
	var result types.Quantity
	if err := c.rpc.Call(ctx, "eth_blobBaseFee", nil, &result); err != nil {
		return nil, err
	}
	return result.BigInt(), nil
}

// GetBalance returns the balance of the given address at the given block.
func (c *Client) GetBalance(ctx context.Context, address types.Address, block BlockNumberOrTag) (*big.Int, error) {
	if block == "" {
		block = BlockLatest
	}

	var result types.Quantity
	if err := c.rpc.Call(ctx, "eth_getBalance", []interface{}{address.String(), block.String()}, &result); err != nil {
		return nil, err
	}
	return result.BigInt(), nil
}

// GetCode returns the code at the given address at the given block.
func (c *Client) GetCode(ctx context.Context, address types.Address, block BlockNumberOrTag) ([]byte, error) {
	if block == "" {
		block = BlockLatest
	}

	var result types.Data
	if err := c.rpc.Call(ctx, "eth_getCode", []interface{}{address.String(), block.String()}, &result); err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}

// GetStorageAt returns the value of a storage slot at the given address.
func (c *Client) GetStorageAt(ctx context.Context, address types.Address, slot types.Hash, block BlockNumberOrTag) (types.Hash, error) {
	if block == "" {
		block = BlockLatest
	}

	var result types.Hash
	if err := c.rpc.Call(ctx, "eth_getStorageAt", []interface{}{address.String(), slot.String(), block.String()}, &result); err != nil {
		return "", err
	}
	return result, nil
}

// GetTransactionCount returns the nonce of the given address at the given block.
func (c *Client) GetTransactionCount(ctx context.Context, address types.Address, block BlockNumberOrTag) (uint64, error) {
	if block == "" {
		block = BlockLatest
	}

	var result types.Quantity
	if err := c.rpc.Call(ctx, "eth_getTransactionCount", []interface{}{address.String(), block.String()}, &result); err != nil {
		return 0, err
	}
	return result.Uint64(), nil
}

// GetBlockByNumber returns a block by its number.
func (c *Client) GetBlockByNumber(ctx context.Context, number BlockNumberOrTag, fullTx bool) (*types.Block, error) {
	if number == "" {
		number = BlockLatest
	}

	var result types.Block
	if err := c.rpc.Call(ctx, "eth_getBlockByNumber", []interface{}{number.String(), fullTx}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetBlockByHash returns a block by its hash.
func (c *Client) GetBlockByHash(ctx context.Context, hash types.Hash, fullTx bool) (*types.Block, error) {
	var result types.Block
	if err := c.rpc.Call(ctx, "eth_getBlockByHash", []interface{}{hash.String(), fullTx}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetBlockTransactionCountByNumber returns the number of transactions in a block by its number.
func (c *Client) GetBlockTransactionCountByNumber(ctx context.Context, number BlockNumberOrTag) (uint64, error) {
	if number == "" {
		number = BlockLatest
	}

	var result types.Quantity
	if err := c.rpc.Call(ctx, "eth_getBlockTransactionCountByNumber", []interface{}{number.String()}, &result); err != nil {
		return 0, err
	}
	return result.Uint64(), nil
}

// GetBlockTransactionCountByHash returns the number of transactions in a block by its hash.
func (c *Client) GetBlockTransactionCountByHash(ctx context.Context, hash types.Hash) (uint64, error) {
	var result types.Quantity
	if err := c.rpc.Call(ctx, "eth_getBlockTransactionCountByHash", []interface{}{hash.String()}, &result); err != nil {
		return 0, err
	}
	return result.Uint64(), nil
}

// GetTransactionByHash returns a transaction by its hash.
func (c *Client) GetTransactionByHash(ctx context.Context, hash types.Hash) (*types.Transaction, error) {
	var result types.Transaction
	if err := c.rpc.Call(ctx, "eth_getTransactionByHash", []interface{}{hash.String()}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTransactionByBlockHashAndIndex returns a transaction by block hash and index.
func (c *Client) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash types.Hash, index uint64) (*types.Transaction, error) {
	var result types.Transaction
	if err := c.rpc.Call(ctx, "eth_getTransactionByBlockHashAndIndex", []interface{}{blockHash.String(), hex.EncodeUint64(index)}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTransactionByBlockNumberAndIndex returns a transaction by block number and index.
func (c *Client) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNumber BlockNumberOrTag, index uint64) (*types.Transaction, error) {
	if blockNumber == "" {
		blockNumber = BlockLatest
	}

	var result types.Transaction
	if err := c.rpc.Call(ctx, "eth_getTransactionByBlockNumberAndIndex", []interface{}{blockNumber.String(), hex.EncodeUint64(index)}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTransactionReceipt returns a transaction receipt by its hash.
func (c *Client) GetTransactionReceipt(ctx context.Context, hash types.Hash) (*types.TransactionReceipt, error) {
	var result types.TransactionReceipt
	if err := c.rpc.Call(ctx, "eth_getTransactionReceipt", []interface{}{hash.String()}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetLogs returns logs matching the given filter.
func (c *Client) GetLogs(ctx context.Context, filter *LogFilter) ([]types.Log, error) {
	var result []types.Log
	if err := c.rpc.Call(ctx, "eth_getLogs", []interface{}{filter}, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Call executes a message call immediately without creating a transaction.
func (c *Client) Call(ctx context.Context, msg *CallMsg, block BlockNumberOrTag) ([]byte, error) {
	if block == "" {
		block = BlockLatest
	}

	var result types.Data
	if err := c.rpc.Call(ctx, "eth_call", []interface{}{msg, block.String()}, &result); err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}

// EstimateGas estimates the gas needed for a transaction.
func (c *Client) EstimateGas(ctx context.Context, msg *CallMsg) (uint64, error) {
	var result types.Quantity
	if err := c.rpc.Call(ctx, "eth_estimateGas", []interface{}{msg}, &result); err != nil {
		return 0, err
	}
	return result.Uint64(), nil
}

// FeeHistory returns historical gas fee data.
func (c *Client) FeeHistory(ctx context.Context, blockCount uint64, newestBlock BlockNumberOrTag, rewardPercentiles []float64) (*FeeHistory, error) {
	if newestBlock == "" {
		newestBlock = BlockLatest
	}

	var result FeeHistory
	if err := c.rpc.Call(ctx, "eth_feeHistory", []interface{}{hex.EncodeUint64(blockCount), newestBlock.String(), rewardPercentiles}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Syncing returns the sync status of the node.
func (c *Client) Syncing(ctx context.Context) (*SyncStatus, error) {
	var result SyncStatus
	if err := c.rpc.Call(ctx, "eth_syncing", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SendRawTransaction sends a signed transaction.
func (c *Client) SendRawTransaction(ctx context.Context, signedTx []byte) (types.Hash, error) {
	var result types.Hash
	if err := c.rpc.Call(ctx, "eth_sendRawTransaction", []interface{}{hex.Encode(signedTx)}, &result); err != nil {
		return "", err
	}
	return result, nil
}

// GetBlockReceipts returns all transaction receipts for a block.
func (c *Client) GetBlockReceipts(ctx context.Context, block BlockNumberOrTag) ([]types.TransactionReceipt, error) {
	if block == "" {
		block = BlockLatest
	}

	var result []types.TransactionReceipt
	if err := c.rpc.Call(ctx, "eth_getBlockReceipts", []interface{}{block.String()}, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetProof returns the account and storage values with Merkle proof.
func (c *Client) GetProof(ctx context.Context, address types.Address, storageKeys []types.Hash, block BlockNumberOrTag) (*AccountProof, error) {
	if block == "" {
		block = BlockLatest
	}

	keys := make([]string, len(storageKeys))
	for i, k := range storageKeys {
		keys[i] = k.String()
	}

	var result AccountProof
	if err := c.rpc.Call(ctx, "eth_getProof", []interface{}{address.String(), keys, block.String()}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AccountProof represents the result of eth_getProof.
type AccountProof struct {
	Address      types.Address  `json:"address"`
	AccountProof []types.Data   `json:"accountProof"`
	Balance      types.Quantity `json:"balance"`
	CodeHash     types.Hash     `json:"codeHash"`
	Nonce        types.Quantity `json:"nonce"`
	StorageHash  types.Hash     `json:"storageHash"`
	StorageProof []StorageProof `json:"storageProof"`
}

// StorageProof represents a storage proof.
type StorageProof struct {
	Key   types.Hash     `json:"key"`
	Value types.Quantity `json:"value"`
	Proof []types.Data   `json:"proof"`
}
