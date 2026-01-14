package node

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ABT-Tech-Limited/alchemy-go/internal/hex"
	"github.com/ABT-Tech-Limited/alchemy-go/types"
)

// BlockNumberOrTag represents a block number or a special tag.
type BlockNumberOrTag string

// Block number tags.
const (
	BlockLatest    BlockNumberOrTag = "latest"
	BlockPending   BlockNumberOrTag = "pending"
	BlockEarliest  BlockNumberOrTag = "earliest"
	BlockFinalized BlockNumberOrTag = "finalized"
	BlockSafe      BlockNumberOrTag = "safe"
)

// BlockNumber creates a BlockNumberOrTag from a block number.
func BlockNumber(n uint64) BlockNumberOrTag {
	return BlockNumberOrTag(hex.EncodeUint64(n))
}

// BlockNumberFromBigInt creates a BlockNumberOrTag from a *big.Int.
func BlockNumberFromBigInt(n *big.Int) BlockNumberOrTag {
	if n == nil {
		return BlockLatest
	}
	return BlockNumberOrTag(hex.EncodeBigInt(n))
}

// String returns the string representation.
func (b BlockNumberOrTag) String() string {
	return string(b)
}

// IsTag returns true if this is a special tag (latest, pending, etc.).
func (b BlockNumberOrTag) IsTag() bool {
	switch b {
	case BlockLatest, BlockPending, BlockEarliest, BlockFinalized, BlockSafe:
		return true
	default:
		return false
	}
}

// Uint64 returns the block number as uint64.
// Returns 0 if this is a tag.
func (b BlockNumberOrTag) Uint64() uint64 {
	if b.IsTag() {
		return 0
	}
	n, _ := hex.DecodeUint64(string(b))
	return n
}

// MarshalJSON implements json.Marshaler.
func (b BlockNumberOrTag) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(b))
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BlockNumberOrTag) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*b = BlockNumberOrTag(s)
	return nil
}

// LogFilter represents a filter for eth_getLogs.
type LogFilter struct {
	// FromBlock is the starting block (inclusive).
	FromBlock BlockNumberOrTag `json:"fromBlock,omitempty"`
	// ToBlock is the ending block (inclusive).
	ToBlock BlockNumberOrTag `json:"toBlock,omitempty"`
	// Address filters logs by contract address.
	Address interface{} `json:"address,omitempty"` // string or []string
	// Topics filters logs by topics.
	Topics []interface{} `json:"topics,omitempty"` // each element is string or []string or nil
	// BlockHash filters logs from a specific block (mutually exclusive with FromBlock/ToBlock).
	BlockHash *types.Hash `json:"blockHash,omitempty"`
}

// NewLogFilter creates a new LogFilter.
func NewLogFilter() *LogFilter {
	return &LogFilter{}
}

// SetFromBlock sets the starting block.
func (f *LogFilter) SetFromBlock(block BlockNumberOrTag) *LogFilter {
	f.FromBlock = block
	return f
}

// SetToBlock sets the ending block.
func (f *LogFilter) SetToBlock(block BlockNumberOrTag) *LogFilter {
	f.ToBlock = block
	return f
}

// SetBlockRange sets both from and to blocks.
func (f *LogFilter) SetBlockRange(from, to BlockNumberOrTag) *LogFilter {
	f.FromBlock = from
	f.ToBlock = to
	return f
}

// SetAddress sets a single address filter.
func (f *LogFilter) SetAddress(address types.Address) *LogFilter {
	f.Address = address.String()
	return f
}

// SetAddresses sets multiple address filters.
func (f *LogFilter) SetAddresses(addresses []types.Address) *LogFilter {
	addrs := make([]string, len(addresses))
	for i, addr := range addresses {
		addrs[i] = addr.String()
	}
	f.Address = addrs
	return f
}

// SetTopic0 sets the first topic (event signature).
func (f *LogFilter) SetTopic0(topic types.Hash) *LogFilter {
	f.ensureTopics(1)
	f.Topics[0] = topic.String()
	return f
}

// SetTopic1 sets the second topic.
func (f *LogFilter) SetTopic1(topic types.Hash) *LogFilter {
	f.ensureTopics(2)
	f.Topics[1] = topic.String()
	return f
}

// SetTopic2 sets the third topic.
func (f *LogFilter) SetTopic2(topic types.Hash) *LogFilter {
	f.ensureTopics(3)
	f.Topics[2] = topic.String()
	return f
}

// SetTopic3 sets the fourth topic.
func (f *LogFilter) SetTopic3(topic types.Hash) *LogFilter {
	f.ensureTopics(4)
	f.Topics[3] = topic.String()
	return f
}

// SetTopic0Or sets the first topic to match any of the given topics.
func (f *LogFilter) SetTopic0Or(topics []types.Hash) *LogFilter {
	f.ensureTopics(1)
	f.Topics[0] = hashesToStrings(topics)
	return f
}

// SetBlockHash sets a specific block hash filter.
func (f *LogFilter) SetBlockHash(hash types.Hash) *LogFilter {
	f.BlockHash = &hash
	f.FromBlock = ""
	f.ToBlock = ""
	return f
}

func (f *LogFilter) ensureTopics(size int) {
	if len(f.Topics) < size {
		newTopics := make([]interface{}, size)
		copy(newTopics, f.Topics)
		f.Topics = newTopics
	}
}

func hashesToStrings(hashes []types.Hash) []string {
	strs := make([]string, len(hashes))
	for i, h := range hashes {
		strs[i] = h.String()
	}
	return strs
}

// CallMsg represents a contract call message.
type CallMsg struct {
	// From is the sender address.
	From *types.Address `json:"from,omitempty"`
	// To is the recipient address (contract address).
	To *types.Address `json:"to,omitempty"`
	// Gas is the gas limit.
	Gas *uint64 `json:"gas,omitempty"`
	// GasPrice is the gas price (legacy).
	GasPrice *big.Int `json:"gasPrice,omitempty"`
	// MaxFeePerGas is the max fee per gas (EIP-1559).
	MaxFeePerGas *big.Int `json:"maxFeePerGas,omitempty"`
	// MaxPriorityFeePerGas is the max priority fee per gas (EIP-1559).
	MaxPriorityFeePerGas *big.Int `json:"maxPriorityFeePerGas,omitempty"`
	// Value is the value to send.
	Value *big.Int `json:"value,omitempty"`
	// Data is the input data.
	Data []byte `json:"data,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (m CallMsg) MarshalJSON() ([]byte, error) {
	type callMsgJSON struct {
		From                 *types.Address `json:"from,omitempty"`
		To                   *types.Address `json:"to,omitempty"`
		Gas                  string         `json:"gas,omitempty"`
		GasPrice             string         `json:"gasPrice,omitempty"`
		MaxFeePerGas         string         `json:"maxFeePerGas,omitempty"`
		MaxPriorityFeePerGas string         `json:"maxPriorityFeePerGas,omitempty"`
		Value                string         `json:"value,omitempty"`
		Data                 string         `json:"data,omitempty"`
	}

	msg := callMsgJSON{
		From: m.From,
		To:   m.To,
	}

	if m.Gas != nil {
		msg.Gas = hex.EncodeUint64(*m.Gas)
	}
	if m.GasPrice != nil {
		msg.GasPrice = hex.EncodeBigInt(m.GasPrice)
	}
	if m.MaxFeePerGas != nil {
		msg.MaxFeePerGas = hex.EncodeBigInt(m.MaxFeePerGas)
	}
	if m.MaxPriorityFeePerGas != nil {
		msg.MaxPriorityFeePerGas = hex.EncodeBigInt(m.MaxPriorityFeePerGas)
	}
	if m.Value != nil {
		msg.Value = hex.EncodeBigInt(m.Value)
	}
	if len(m.Data) > 0 {
		msg.Data = hex.Encode(m.Data)
	}

	return json.Marshal(msg)
}

// FeeHistory represents the result of eth_feeHistory.
type FeeHistory struct {
	// OldestBlock is the oldest block in the range.
	OldestBlock types.Quantity `json:"oldestBlock"`
	// BaseFeePerGas is the base fee per gas for each block.
	BaseFeePerGas []types.Quantity `json:"baseFeePerGas"`
	// GasUsedRatio is the ratio of gas used to gas limit for each block.
	GasUsedRatio []float64 `json:"gasUsedRatio"`
	// Reward is the priority fee per gas percentiles for each block.
	Reward [][]types.Quantity `json:"reward,omitempty"`
	// BaseFeePerBlobGas is the base fee per blob gas for each block (EIP-4844).
	BaseFeePerBlobGas []types.Quantity `json:"baseFeePerBlobGas,omitempty"`
	// BlobGasUsedRatio is the ratio of blob gas used for each block (EIP-4844).
	BlobGasUsedRatio []float64 `json:"blobGasUsedRatio,omitempty"`
}

// SyncStatus represents the result of eth_syncing.
type SyncStatus struct {
	// Syncing is true if the node is syncing.
	Syncing bool
	// StartingBlock is the block at which the sync started.
	StartingBlock uint64
	// CurrentBlock is the current block being synced.
	CurrentBlock uint64
	// HighestBlock is the highest known block.
	HighestBlock uint64
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SyncStatus) UnmarshalJSON(data []byte) error {
	// First try to unmarshal as boolean (not syncing)
	var syncing bool
	if err := json.Unmarshal(data, &syncing); err == nil {
		s.Syncing = syncing
		return nil
	}

	// Then try to unmarshal as sync status object
	type syncStatusJSON struct {
		StartingBlock types.Quantity `json:"startingBlock"`
		CurrentBlock  types.Quantity `json:"currentBlock"`
		HighestBlock  types.Quantity `json:"highestBlock"`
	}

	var status syncStatusJSON
	if err := json.Unmarshal(data, &status); err != nil {
		return fmt.Errorf("failed to unmarshal sync status: %w", err)
	}

	s.Syncing = true
	s.StartingBlock = status.StartingBlock.Uint64()
	s.CurrentBlock = status.CurrentBlock.Uint64()
	s.HighestBlock = status.HighestBlock.Uint64()

	return nil
}

// TraceConfig represents configuration for trace methods.
type TraceConfig struct {
	// DisableStorage disables storage capture.
	DisableStorage bool `json:"disableStorage,omitempty"`
	// DisableStack disables stack capture.
	DisableStack bool `json:"disableStack,omitempty"`
	// EnableMemory enables memory capture.
	EnableMemory bool `json:"enableMemory,omitempty"`
	// EnableReturnData enables return data capture.
	EnableReturnData bool `json:"enableReturnData,omitempty"`
	// Tracer is the name of a built-in tracer.
	Tracer string `json:"tracer,omitempty"`
	// TracerConfig is configuration for the tracer.
	TracerConfig json.RawMessage `json:"tracerConfig,omitempty"`
	// Timeout is the maximum time for tracing.
	Timeout string `json:"timeout,omitempty"`
}
