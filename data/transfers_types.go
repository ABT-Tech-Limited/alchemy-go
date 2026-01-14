package data

import (
	"github.com/ABT-Tech-Limited/alchemy-go/types"
)

// AssetTransferCategory represents the type of asset transfer.
type AssetTransferCategory string

// Asset transfer categories.
const (
	CategoryExternal   AssetTransferCategory = "external"
	CategoryInternal   AssetTransferCategory = "internal"
	CategoryERC20      AssetTransferCategory = "erc20"
	CategoryERC721     AssetTransferCategory = "erc721"
	CategoryERC1155    AssetTransferCategory = "erc1155"
	CategorySpecialNFT AssetTransferCategory = "specialnft"
)

// SortOrder represents the sort order for results.
type SortOrder string

// Sort orders.
const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

// AssetTransfersParams represents the parameters for getAssetTransfers.
type AssetTransfersParams struct {
	// FromBlock is the starting block (hex or "latest").
	FromBlock string `json:"fromBlock,omitempty"`
	// ToBlock is the ending block (hex or "latest").
	ToBlock string `json:"toBlock,omitempty"`
	// FromAddress filters transfers from this address.
	FromAddress *types.Address `json:"fromAddress,omitempty"`
	// ToAddress filters transfers to this address.
	ToAddress *types.Address `json:"toAddress,omitempty"`
	// ContractAddresses filters transfers by contract addresses.
	ContractAddresses []types.Address `json:"contractAddresses,omitempty"`
	// Category is the list of transfer categories to include.
	Category []AssetTransferCategory `json:"category"`
	// Order is the sort order (asc or desc).
	Order SortOrder `json:"order,omitempty"`
	// WithMetadata includes block timestamps in the response.
	WithMetadata bool `json:"withMetadata,omitempty"`
	// ExcludeZeroValue excludes zero-value transfers.
	ExcludeZeroValue bool `json:"excludeZeroValue,omitempty"`
	// MaxCount is the maximum number of results per page (hex).
	MaxCount string `json:"maxCount,omitempty"`
	// PageKey is the pagination key for fetching more results.
	PageKey string `json:"pageKey,omitempty"`
}

// NewAssetTransfersParams creates a new AssetTransfersParams with default values.
func NewAssetTransfersParams() *AssetTransfersParams {
	return &AssetTransfersParams{
		Category:         []AssetTransferCategory{CategoryExternal, CategoryERC20},
		ExcludeZeroValue: true,
	}
}

// SetFromBlock sets the starting block.
func (p *AssetTransfersParams) SetFromBlock(block string) *AssetTransfersParams {
	p.FromBlock = block
	return p
}

// SetToBlock sets the ending block.
func (p *AssetTransfersParams) SetToBlock(block string) *AssetTransfersParams {
	p.ToBlock = block
	return p
}

// SetFromAddress sets the from address filter.
func (p *AssetTransfersParams) SetFromAddress(address types.Address) *AssetTransfersParams {
	p.FromAddress = &address
	return p
}

// SetToAddress sets the to address filter.
func (p *AssetTransfersParams) SetToAddress(address types.Address) *AssetTransfersParams {
	p.ToAddress = &address
	return p
}

// SetContractAddresses sets the contract address filter.
func (p *AssetTransfersParams) SetContractAddresses(addresses []types.Address) *AssetTransfersParams {
	p.ContractAddresses = addresses
	return p
}

// SetCategories sets the transfer categories.
func (p *AssetTransfersParams) SetCategories(categories []AssetTransferCategory) *AssetTransfersParams {
	p.Category = categories
	return p
}

// SetOrder sets the sort order.
func (p *AssetTransfersParams) SetOrder(order SortOrder) *AssetTransfersParams {
	p.Order = order
	return p
}

// SetWithMetadata enables metadata in the response.
func (p *AssetTransfersParams) SetWithMetadata(withMetadata bool) *AssetTransfersParams {
	p.WithMetadata = withMetadata
	return p
}

// SetMaxCount sets the maximum number of results.
func (p *AssetTransfersParams) SetMaxCount(count int) *AssetTransfersParams {
	p.MaxCount = "0x" + string(rune(count))
	return p
}

// AssetTransfersResponse represents the response from getAssetTransfers.
type AssetTransfersResponse struct {
	// PageKey is the pagination key for fetching more results.
	PageKey string `json:"pageKey,omitempty"`
	// Transfers is the list of asset transfers.
	Transfers []AssetTransfer `json:"transfers"`
}

// HasMore returns true if there are more results available.
func (r *AssetTransfersResponse) HasMore() bool {
	return r.PageKey != ""
}

// AssetTransfer represents a single asset transfer.
type AssetTransfer struct {
	// Category is the type of transfer.
	Category AssetTransferCategory `json:"category"`
	// BlockNum is the block number (hex).
	BlockNum string `json:"blockNum"`
	// From is the sender address.
	From types.Address `json:"from"`
	// To is the recipient address.
	To *types.Address `json:"to,omitempty"`
	// Value is the transferred value (as decimal number).
	Value *float64 `json:"value,omitempty"`
	// TokenID is the NFT token ID.
	TokenID *string `json:"tokenId,omitempty"`
	// ERC1155Metadata contains ERC1155 token metadata.
	ERC1155Metadata []ERC1155Metadata `json:"erc1155Metadata,omitempty"`
	// Asset is the asset symbol (e.g., "ETH", "USDC").
	Asset *string `json:"asset,omitempty"`
	// UniqueID is the unique identifier for this transfer.
	UniqueID string `json:"uniqueId"`
	// Hash is the transaction hash.
	Hash types.Hash `json:"hash"`
	// RawContract contains raw contract information.
	RawContract RawContract `json:"rawContract"`
	// Metadata contains additional metadata (when WithMetadata is true).
	Metadata *TransferMetadata `json:"metadata,omitempty"`
}

// BlockNumber returns the block number as uint64.
func (t *AssetTransfer) BlockNumber() uint64 {
	if t.BlockNum == "" {
		return 0
	}
	// Parse hex block number
	var n uint64
	for _, c := range t.BlockNum[2:] {
		n = n*16 + uint64(hexDigitToValue(byte(c)))
	}
	return n
}

func hexDigitToValue(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c - 'a' + 10)
	case c >= 'A' && c <= 'F':
		return int(c - 'A' + 10)
	}
	return 0
}

// RawContract contains raw contract information.
type RawContract struct {
	// Value is the raw transfer value (hex).
	Value *string `json:"value,omitempty"`
	// Address is the contract address.
	Address *string `json:"address,omitempty"`
	// Decimal is the token decimals (hex).
	Decimal *string `json:"decimal,omitempty"`
}

// ERC1155Metadata contains ERC1155 token metadata.
type ERC1155Metadata struct {
	// TokenID is the token ID (hex).
	TokenID string `json:"tokenId"`
	// Value is the amount transferred (hex).
	Value string `json:"value"`
}

// TransferMetadata contains additional transfer metadata.
type TransferMetadata struct {
	// BlockTimestamp is the block timestamp in ISO format.
	BlockTimestamp string `json:"blockTimestamp"`
}
