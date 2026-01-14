package data

import (
	"github.com/ABT-Tech-Limited/alchemy-go/types"
)

// NFTFilter represents a filter for NFT queries.
type NFTFilter string

// NFT filters.
const (
	NFTFilterSpam     NFTFilter = "SPAM"
	NFTFilterAirdrops NFTFilter = "AIRDROPS"
)

// SpamConfidenceLevel represents the spam confidence level.
type SpamConfidenceLevel string

// Spam confidence levels.
const (
	SpamConfidenceVeryHigh SpamConfidenceLevel = "VERY_HIGH"
	SpamConfidenceHigh     SpamConfidenceLevel = "HIGH"
	SpamConfidenceMedium   SpamConfidenceLevel = "MEDIUM"
	SpamConfidenceLow      SpamConfidenceLevel = "LOW"
)

// NFTOrderBy represents the ordering for NFT results.
type NFTOrderBy string

// NFT order options.
const (
	NFTOrderByTransferTime NFTOrderBy = "transferTime"
)

// NFTTokenType represents the type of NFT token.
type NFTTokenType string

// NFT token types.
const (
	NFTTokenTypeERC721  NFTTokenType = "ERC721"
	NFTTokenTypeERC1155 NFTTokenType = "ERC1155"
)

// NFTsForOwnerParams represents the parameters for getNFTsForOwner.
type NFTsForOwnerParams struct {
	// Owner is the address of the NFT owner.
	Owner types.Address `json:"owner"`
	// ContractAddresses filters NFTs by contract addresses (max 45).
	ContractAddresses []types.Address `json:"contractAddresses,omitempty"`
	// WithMetadata includes NFT metadata in the response.
	WithMetadata *bool `json:"withMetadata,omitempty"`
	// OrderBy specifies the ordering of results.
	OrderBy NFTOrderBy `json:"orderBy,omitempty"`
	// ExcludeFilters excludes NFTs matching these filters.
	ExcludeFilters []NFTFilter `json:"excludeFilters,omitempty"`
	// IncludeFilters includes only NFTs matching these filters.
	IncludeFilters []NFTFilter `json:"includeFilters,omitempty"`
	// SpamConfidenceLevel sets the spam detection threshold.
	SpamConfidenceLevel SpamConfidenceLevel `json:"spamConfidenceLevel,omitempty"`
	// TokenURITimeoutInMs sets the timeout for fetching token URIs.
	TokenURITimeoutInMs *int `json:"tokenUriTimeoutInMs,omitempty"`
	// PageKey is the pagination key.
	PageKey string `json:"pageKey,omitempty"`
	// PageSize is the number of results per page (max 100).
	PageSize *int `json:"pageSize,omitempty"`
}

// NewNFTsForOwnerParams creates new NFTsForOwnerParams.
func NewNFTsForOwnerParams(owner types.Address) *NFTsForOwnerParams {
	return &NFTsForOwnerParams{
		Owner: owner,
	}
}

// SetContractAddresses sets the contract address filter.
func (p *NFTsForOwnerParams) SetContractAddresses(addresses []types.Address) *NFTsForOwnerParams {
	p.ContractAddresses = addresses
	return p
}

// SetWithMetadata enables metadata in the response.
func (p *NFTsForOwnerParams) SetWithMetadata(withMetadata bool) *NFTsForOwnerParams {
	p.WithMetadata = &withMetadata
	return p
}

// SetOrderBy sets the ordering.
func (p *NFTsForOwnerParams) SetOrderBy(orderBy NFTOrderBy) *NFTsForOwnerParams {
	p.OrderBy = orderBy
	return p
}

// SetExcludeFilters sets exclusion filters.
func (p *NFTsForOwnerParams) SetExcludeFilters(filters []NFTFilter) *NFTsForOwnerParams {
	p.ExcludeFilters = filters
	return p
}

// SetPageSize sets the page size.
func (p *NFTsForOwnerParams) SetPageSize(size int) *NFTsForOwnerParams {
	p.PageSize = &size
	return p
}

// NFTsForOwnerResponse represents the response from getNFTsForOwner.
type NFTsForOwnerResponse struct {
	// OwnedNFTs is the list of owned NFTs.
	OwnedNFTs []OwnedNFT `json:"ownedNfts"`
	// TotalCount is the total number of NFTs.
	TotalCount int `json:"totalCount"`
	// PageKey is the pagination key.
	PageKey string `json:"pageKey,omitempty"`
	// ValidAt contains block information.
	ValidAt *ValidAt `json:"validAt,omitempty"`
}

// HasMore returns true if there are more results available.
func (r *NFTsForOwnerResponse) HasMore() bool {
	return r.PageKey != ""
}

// OwnedNFT represents an NFT owned by an address.
type OwnedNFT struct {
	// Contract contains contract information.
	Contract NFTContract `json:"contract"`
	// TokenID is the token ID.
	TokenID string `json:"tokenId"`
	// TokenType is the token type (ERC721, ERC1155).
	TokenType string `json:"tokenType"`
	// Name is the NFT name.
	Name *string `json:"name,omitempty"`
	// Description is the NFT description.
	Description *string `json:"description,omitempty"`
	// Image contains image information.
	Image *NFTImage `json:"image,omitempty"`
	// Raw contains raw token data.
	Raw *NFTRaw `json:"raw,omitempty"`
	// Collection contains collection information.
	Collection *NFTCollection `json:"collection,omitempty"`
	// TokenURI is the token URI.
	TokenURI *string `json:"tokenUri,omitempty"`
	// TimeLastUpdated is when the metadata was last updated.
	TimeLastUpdated *string `json:"timeLastUpdated,omitempty"`
	// AcquiredAt contains acquisition information.
	AcquiredAt *AcquiredAt `json:"acquiredAt,omitempty"`
	// Balance is the token balance (for ERC1155).
	Balance *string `json:"balance,omitempty"`
}

// NFTContract represents NFT contract information.
type NFTContract struct {
	// Address is the contract address.
	Address types.Address `json:"address"`
	// Name is the contract name.
	Name *string `json:"name,omitempty"`
	// Symbol is the contract symbol.
	Symbol *string `json:"symbol,omitempty"`
	// TotalSupply is the total supply.
	TotalSupply *string `json:"totalSupply,omitempty"`
	// TokenType is the token type.
	TokenType string `json:"tokenType"`
	// ContractDeployer is the deployer address.
	ContractDeployer *string `json:"contractDeployer,omitempty"`
	// DeployedBlockNumber is the deployment block number.
	DeployedBlockNumber *int64 `json:"deployedBlockNumber,omitempty"`
	// OpenSeaMetadata contains OpenSea metadata.
	OpenSeaMetadata *OpenSeaMetadata `json:"openseaMetadata,omitempty"`
	// IsSpam indicates if the contract is marked as spam.
	IsSpam *bool `json:"isSpam,omitempty"`
	// SpamClassifications contains spam classification reasons.
	SpamClassifications []string `json:"spamClassifications,omitempty"`
}

// OpenSeaMetadata contains OpenSea-specific metadata.
type OpenSeaMetadata struct {
	// FloorPrice is the floor price.
	FloorPrice *float64 `json:"floorPrice,omitempty"`
	// CollectionName is the collection name.
	CollectionName *string `json:"collectionName,omitempty"`
	// SafelistRequestStatus is the safelist status.
	SafelistRequestStatus *string `json:"safelistRequestStatus,omitempty"`
	// ImageURL is the collection image URL.
	ImageURL *string `json:"imageUrl,omitempty"`
	// Description is the collection description.
	Description *string `json:"description,omitempty"`
	// ExternalURL is the external URL.
	ExternalURL *string `json:"externalUrl,omitempty"`
	// TwitterUsername is the Twitter username.
	TwitterUsername *string `json:"twitterUsername,omitempty"`
	// DiscordURL is the Discord URL.
	DiscordURL *string `json:"discordUrl,omitempty"`
	// BannerImageURL is the banner image URL.
	BannerImageURL *string `json:"bannerImageUrl,omitempty"`
	// LastIngestedAt is when the metadata was last ingested.
	LastIngestedAt *string `json:"lastIngestedAt,omitempty"`
}

// NFTImage contains NFT image information.
type NFTImage struct {
	// CachedURL is the cached image URL.
	CachedURL *string `json:"cachedUrl,omitempty"`
	// ThumbnailURL is the thumbnail URL.
	ThumbnailURL *string `json:"thumbnailUrl,omitempty"`
	// PngURL is the PNG URL.
	PngURL *string `json:"pngUrl,omitempty"`
	// ContentType is the content type.
	ContentType *string `json:"contentType,omitempty"`
	// Size is the image size in bytes.
	Size *int `json:"size,omitempty"`
	// OriginalURL is the original image URL.
	OriginalURL *string `json:"originalUrl,omitempty"`
}

// NFTRaw contains raw NFT data.
type NFTRaw struct {
	// TokenURI is the raw token URI.
	TokenURI *string `json:"tokenUri,omitempty"`
	// Metadata contains raw metadata.
	Metadata *NFTRawMetadata `json:"metadata,omitempty"`
	// Error is any error during metadata fetch.
	Error *string `json:"error,omitempty"`
}

// NFTRawMetadata contains raw NFT metadata.
type NFTRawMetadata struct {
	// Image is the image URL.
	Image *string `json:"image,omitempty"`
	// Name is the NFT name.
	Name *string `json:"name,omitempty"`
	// Description is the NFT description.
	Description *string `json:"description,omitempty"`
	// Attributes is the list of attributes.
	Attributes []NFTAttribute `json:"attributes,omitempty"`
	// ExternalURL is the external URL.
	ExternalURL *string `json:"external_url,omitempty"`
	// AnimationURL is the animation URL.
	AnimationURL *string `json:"animation_url,omitempty"`
}

// NFTAttribute represents an NFT attribute.
type NFTAttribute struct {
	// TraitType is the trait type.
	TraitType *string `json:"trait_type,omitempty"`
	// Value is the trait value.
	Value interface{} `json:"value,omitempty"`
	// DisplayType is the display type.
	DisplayType *string `json:"display_type,omitempty"`
}

// NFTCollection contains collection information.
type NFTCollection struct {
	// Name is the collection name.
	Name *string `json:"name,omitempty"`
	// Slug is the collection slug.
	Slug *string `json:"slug,omitempty"`
	// ExternalURL is the external URL.
	ExternalURL *string `json:"externalUrl,omitempty"`
	// BannerImageURL is the banner image URL.
	BannerImageURL *string `json:"bannerImageUrl,omitempty"`
}

// ValidAt contains block validity information.
type ValidAt struct {
	// BlockNumber is the block number.
	BlockNumber int `json:"blockNumber"`
	// BlockHash is the block hash.
	BlockHash string `json:"blockHash"`
	// BlockTimestamp is the block timestamp.
	BlockTimestamp string `json:"blockTimestamp"`
}

// AcquiredAt contains NFT acquisition information.
type AcquiredAt struct {
	// BlockTimestamp is when the NFT was acquired.
	BlockTimestamp *string `json:"blockTimestamp,omitempty"`
	// BlockNumber is the block number.
	BlockNumber *string `json:"blockNumber,omitempty"`
}

// NFTMetadataParams represents the parameters for getNFTMetadata.
type NFTMetadataParams struct {
	// ContractAddress is the NFT contract address.
	ContractAddress types.Address `json:"contractAddress"`
	// TokenID is the token ID.
	TokenID string `json:"tokenId"`
	// TokenType is the token type (optional).
	TokenType *string `json:"tokenType,omitempty"`
	// RefreshCache forces a metadata refresh.
	RefreshCache bool `json:"refreshCache,omitempty"`
}

// NewNFTMetadataParams creates new NFTMetadataParams.
func NewNFTMetadataParams(contractAddress types.Address, tokenID string) *NFTMetadataParams {
	return &NFTMetadataParams{
		ContractAddress: contractAddress,
		TokenID:         tokenID,
	}
}

// SetTokenType sets the token type.
func (p *NFTMetadataParams) SetTokenType(tokenType string) *NFTMetadataParams {
	p.TokenType = &tokenType
	return p
}

// SetRefreshCache enables cache refresh.
func (p *NFTMetadataParams) SetRefreshCache(refresh bool) *NFTMetadataParams {
	p.RefreshCache = refresh
	return p
}

// NFTContractMetadata represents contract-level NFT metadata.
type NFTContractMetadata struct {
	// Address is the contract address.
	Address types.Address `json:"address"`
	// Name is the contract name.
	Name *string `json:"name,omitempty"`
	// Symbol is the contract symbol.
	Symbol *string `json:"symbol,omitempty"`
	// TotalSupply is the total supply.
	TotalSupply *string `json:"totalSupply,omitempty"`
	// TokenType is the token type.
	TokenType string `json:"tokenType"`
	// ContractDeployer is the deployer address.
	ContractDeployer *string `json:"contractDeployer,omitempty"`
	// DeployedBlockNumber is the deployment block number.
	DeployedBlockNumber *int64 `json:"deployedBlockNumber,omitempty"`
	// OpenSeaMetadata contains OpenSea metadata.
	OpenSeaMetadata *OpenSeaMetadata `json:"openseaMetadata,omitempty"`
}
