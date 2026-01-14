package types

// TransactionReceipt represents an Ethereum transaction receipt.
type TransactionReceipt struct {
	// TransactionHash is the transaction hash.
	TransactionHash Hash `json:"transactionHash"`
	// TransactionIndex is the transaction index in the block.
	TransactionIndex Quantity `json:"transactionIndex"`
	// BlockHash is the block hash.
	BlockHash Hash `json:"blockHash"`
	// BlockNumber is the block number.
	BlockNumber Quantity `json:"blockNumber"`
	// From is the sender address.
	From Address `json:"from"`
	// To is the recipient address (null for contract creation).
	To *Address `json:"to,omitempty"`
	// CumulativeGasUsed is the cumulative gas used in the block up to this transaction.
	CumulativeGasUsed Quantity `json:"cumulativeGasUsed"`
	// GasUsed is the gas used by this transaction.
	GasUsed Quantity `json:"gasUsed"`
	// EffectiveGasPrice is the effective gas price (EIP-1559).
	EffectiveGasPrice Quantity `json:"effectiveGasPrice"`
	// ContractAddress is the contract address created (if any).
	ContractAddress *Address `json:"contractAddress,omitempty"`
	// Logs is the list of logs emitted by the transaction.
	Logs []Log `json:"logs"`
	// LogsBloom is the bloom filter for logs.
	LogsBloom Data `json:"logsBloom"`
	// Type is the transaction type.
	Type Quantity `json:"type"`
	// Status is the transaction status (1 = success, 0 = failure).
	Status Quantity `json:"status"`
	// Root is the post-transaction state root (pre-Byzantium only).
	Root *Hash `json:"root,omitempty"`
	// BlobGasUsed is the blob gas used (EIP-4844).
	BlobGasUsed *Quantity `json:"blobGasUsed,omitempty"`
	// BlobGasPrice is the blob gas price (EIP-4844).
	BlobGasPrice *Quantity `json:"blobGasPrice,omitempty"`
}

// IsSuccessful returns true if the transaction was successful.
func (r *TransactionReceipt) IsSuccessful() bool {
	return r.Status.Uint64() == 1
}

// IsFailed returns true if the transaction failed.
func (r *TransactionReceipt) IsFailed() bool {
	return r.Status.Uint64() == 0
}

// IsContractCreation returns true if the transaction created a contract.
func (r *TransactionReceipt) IsContractCreation() bool {
	return r.ContractAddress != nil && !r.ContractAddress.IsZero()
}
