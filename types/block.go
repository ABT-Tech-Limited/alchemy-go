package types

import (
	"encoding/json"
)

// Block represents an Ethereum block.
type Block struct {
	// Hash is the block hash.
	Hash Hash `json:"hash"`
	// ParentHash is the parent block hash.
	ParentHash Hash `json:"parentHash"`
	// Sha3Uncles is the SHA3 of the uncles data in the block.
	Sha3Uncles Hash `json:"sha3Uncles"`
	// Miner is the address of the miner.
	Miner Address `json:"miner"`
	// StateRoot is the root of the state trie.
	StateRoot Hash `json:"stateRoot"`
	// TransactionsRoot is the root of the transaction trie.
	TransactionsRoot Hash `json:"transactionsRoot"`
	// ReceiptsRoot is the root of the receipts trie.
	ReceiptsRoot Hash `json:"receiptsRoot"`
	// LogsBloom is the bloom filter for logs.
	LogsBloom Data `json:"logsBloom"`
	// Difficulty is the difficulty of the block.
	Difficulty Quantity `json:"difficulty,omitempty"`
	// TotalDifficulty is the total difficulty of the chain until this block.
	TotalDifficulty Quantity `json:"totalDifficulty,omitempty"`
	// Number is the block number.
	Number Quantity `json:"number"`
	// GasLimit is the gas limit for the block.
	GasLimit Quantity `json:"gasLimit"`
	// GasUsed is the gas used by all transactions in the block.
	GasUsed Quantity `json:"gasUsed"`
	// Timestamp is the block timestamp.
	Timestamp Quantity `json:"timestamp"`
	// ExtraData is the extra data field of the block.
	ExtraData Data `json:"extraData"`
	// MixHash is the mix hash.
	MixHash Hash `json:"mixHash"`
	// Nonce is the block nonce.
	Nonce Data `json:"nonce"`
	// BaseFeePerGas is the base fee per gas (EIP-1559).
	BaseFeePerGas *Quantity `json:"baseFeePerGas,omitempty"`
	// WithdrawalsRoot is the root of the withdrawals trie (Shanghai upgrade).
	WithdrawalsRoot *Hash `json:"withdrawalsRoot,omitempty"`
	// BlobGasUsed is the blob gas used (EIP-4844).
	BlobGasUsed *Quantity `json:"blobGasUsed,omitempty"`
	// ExcessBlobGas is the excess blob gas (EIP-4844).
	ExcessBlobGas *Quantity `json:"excessBlobGas,omitempty"`
	// ParentBeaconBlockRoot is the parent beacon block root.
	ParentBeaconBlockRoot *Hash `json:"parentBeaconBlockRoot,omitempty"`
	// Size is the block size in bytes.
	Size Quantity `json:"size"`
	// Uncles is the list of uncle block hashes.
	Uncles []Hash `json:"uncles"`
	// Withdrawals is the list of withdrawals (Shanghai upgrade).
	Withdrawals []Withdrawal `json:"withdrawals,omitempty"`

	// Transactions can be either a list of transaction hashes or full transaction objects.
	// Use TransactionHashes() or Transactions() to access them.
	rawTransactions json.RawMessage
}

// rawBlock is used for unmarshaling to handle transactions field.
type rawBlock struct {
	Hash                  Hash            `json:"hash"`
	ParentHash            Hash            `json:"parentHash"`
	Sha3Uncles            Hash            `json:"sha3Uncles"`
	Miner                 Address         `json:"miner"`
	StateRoot             Hash            `json:"stateRoot"`
	TransactionsRoot      Hash            `json:"transactionsRoot"`
	ReceiptsRoot          Hash            `json:"receiptsRoot"`
	LogsBloom             Data            `json:"logsBloom"`
	Difficulty            Quantity        `json:"difficulty,omitempty"`
	TotalDifficulty       Quantity        `json:"totalDifficulty,omitempty"`
	Number                Quantity        `json:"number"`
	GasLimit              Quantity        `json:"gasLimit"`
	GasUsed               Quantity        `json:"gasUsed"`
	Timestamp             Quantity        `json:"timestamp"`
	ExtraData             Data            `json:"extraData"`
	MixHash               Hash            `json:"mixHash"`
	Nonce                 Data            `json:"nonce"`
	BaseFeePerGas         *Quantity       `json:"baseFeePerGas,omitempty"`
	WithdrawalsRoot       *Hash           `json:"withdrawalsRoot,omitempty"`
	BlobGasUsed           *Quantity       `json:"blobGasUsed,omitempty"`
	ExcessBlobGas         *Quantity       `json:"excessBlobGas,omitempty"`
	ParentBeaconBlockRoot *Hash           `json:"parentBeaconBlockRoot,omitempty"`
	Size                  Quantity        `json:"size"`
	Uncles                []Hash          `json:"uncles"`
	Withdrawals           []Withdrawal    `json:"withdrawals,omitempty"`
	Transactions          json.RawMessage `json:"transactions"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Block) UnmarshalJSON(data []byte) error {
	var raw rawBlock
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	b.Hash = raw.Hash
	b.ParentHash = raw.ParentHash
	b.Sha3Uncles = raw.Sha3Uncles
	b.Miner = raw.Miner
	b.StateRoot = raw.StateRoot
	b.TransactionsRoot = raw.TransactionsRoot
	b.ReceiptsRoot = raw.ReceiptsRoot
	b.LogsBloom = raw.LogsBloom
	b.Difficulty = raw.Difficulty
	b.TotalDifficulty = raw.TotalDifficulty
	b.Number = raw.Number
	b.GasLimit = raw.GasLimit
	b.GasUsed = raw.GasUsed
	b.Timestamp = raw.Timestamp
	b.ExtraData = raw.ExtraData
	b.MixHash = raw.MixHash
	b.Nonce = raw.Nonce
	b.BaseFeePerGas = raw.BaseFeePerGas
	b.WithdrawalsRoot = raw.WithdrawalsRoot
	b.BlobGasUsed = raw.BlobGasUsed
	b.ExcessBlobGas = raw.ExcessBlobGas
	b.ParentBeaconBlockRoot = raw.ParentBeaconBlockRoot
	b.Size = raw.Size
	b.Uncles = raw.Uncles
	b.Withdrawals = raw.Withdrawals
	b.rawTransactions = raw.Transactions

	return nil
}

// TransactionHashes returns the transaction hashes if the block was fetched without full transactions.
// Returns nil if the block contains full transaction objects.
func (b *Block) TransactionHashes() []Hash {
	if len(b.rawTransactions) == 0 {
		return nil
	}

	var hashes []Hash
	if err := json.Unmarshal(b.rawTransactions, &hashes); err != nil {
		return nil
	}
	return hashes
}

// Transactions returns the full transactions if the block was fetched with full transactions.
// Returns nil if the block only contains transaction hashes.
func (b *Block) Transactions() []Transaction {
	if len(b.rawTransactions) == 0 {
		return nil
	}

	var txs []Transaction
	if err := json.Unmarshal(b.rawTransactions, &txs); err != nil {
		return nil
	}
	return txs
}

// TransactionCount returns the number of transactions in the block.
func (b *Block) TransactionCount() int {
	if len(b.rawTransactions) == 0 {
		return 0
	}

	// Try to count as array elements
	var arr []json.RawMessage
	if err := json.Unmarshal(b.rawTransactions, &arr); err != nil {
		return 0
	}
	return len(arr)
}

// Withdrawal represents a validator withdrawal (Shanghai upgrade).
type Withdrawal struct {
	// Index is the withdrawal index.
	Index Quantity `json:"index"`
	// ValidatorIndex is the validator index.
	ValidatorIndex Quantity `json:"validatorIndex"`
	// Address is the recipient address.
	Address Address `json:"address"`
	// Amount is the amount in Gwei.
	Amount Quantity `json:"amount"`
}

// Transaction represents an Ethereum transaction.
type Transaction struct {
	// Hash is the transaction hash.
	Hash Hash `json:"hash"`
	// Nonce is the sender's nonce.
	Nonce Quantity `json:"nonce"`
	// BlockHash is the block hash (null for pending transactions).
	BlockHash *Hash `json:"blockHash,omitempty"`
	// BlockNumber is the block number (null for pending transactions).
	BlockNumber *Quantity `json:"blockNumber,omitempty"`
	// TransactionIndex is the transaction index in the block (null for pending).
	TransactionIndex *Quantity `json:"transactionIndex,omitempty"`
	// From is the sender address.
	From Address `json:"from"`
	// To is the recipient address (null for contract creation).
	To *Address `json:"to,omitempty"`
	// Value is the value transferred in wei.
	Value Quantity `json:"value"`
	// Gas is the gas limit.
	Gas Quantity `json:"gas"`
	// GasPrice is the gas price (legacy and EIP-2930 transactions).
	GasPrice *Quantity `json:"gasPrice,omitempty"`
	// MaxFeePerGas is the max fee per gas (EIP-1559 transactions).
	MaxFeePerGas *Quantity `json:"maxFeePerGas,omitempty"`
	// MaxPriorityFeePerGas is the max priority fee per gas (EIP-1559 transactions).
	MaxPriorityFeePerGas *Quantity `json:"maxPriorityFeePerGas,omitempty"`
	// MaxFeePerBlobGas is the max fee per blob gas (EIP-4844 transactions).
	MaxFeePerBlobGas *Quantity `json:"maxFeePerBlobGas,omitempty"`
	// Input is the transaction input data.
	Input Data `json:"input"`
	// V is the recovery id.
	V Quantity `json:"v"`
	// R is the signature r value.
	R Quantity `json:"r"`
	// S is the signature s value.
	S Quantity `json:"s"`
	// YParity is the y-parity of the signature (EIP-1559 and later).
	YParity *Quantity `json:"yParity,omitempty"`
	// Type is the transaction type.
	Type *Quantity `json:"type,omitempty"`
	// ChainID is the chain ID.
	ChainID *Quantity `json:"chainId,omitempty"`
	// AccessList is the access list (EIP-2930 and later).
	AccessList []AccessListEntry `json:"accessList,omitempty"`
	// BlobVersionedHashes is the blob versioned hashes (EIP-4844).
	BlobVersionedHashes []Hash `json:"blobVersionedHashes,omitempty"`
}

// AccessListEntry represents an entry in an access list.
type AccessListEntry struct {
	// Address is the accessed address.
	Address Address `json:"address"`
	// StorageKeys is the list of accessed storage keys.
	StorageKeys []Hash `json:"storageKeys"`
}

// TxType returns the transaction type as an integer.
func (tx *Transaction) TxType() int {
	if tx.Type == nil {
		return 0
	}
	return int(tx.Type.Uint64())
}

// IsLegacy returns true if this is a legacy (type 0) transaction.
func (tx *Transaction) IsLegacy() bool {
	return tx.TxType() == 0
}

// IsAccessList returns true if this is an access list (type 1) transaction.
func (tx *Transaction) IsAccessList() bool {
	return tx.TxType() == 1
}

// IsDynamicFee returns true if this is a dynamic fee (type 2) transaction.
func (tx *Transaction) IsDynamicFee() bool {
	return tx.TxType() == 2
}

// IsBlob returns true if this is a blob (type 3) transaction.
func (tx *Transaction) IsBlob() bool {
	return tx.TxType() == 3
}
