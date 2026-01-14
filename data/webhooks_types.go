package data

// WebhookType represents the type of webhook.
type WebhookType string

// Webhook types.
const (
	WebhookTypeGraphQL         WebhookType = "GRAPHQL"
	WebhookTypeAddressActivity WebhookType = "ADDRESS_ACTIVITY"
	WebhookTypeNFTActivity     WebhookType = "NFT_ACTIVITY"
)

// WebhookNetwork represents the network for a webhook.
type WebhookNetwork string

// Webhook networks.
const (
	// Ethereum
	WebhookNetworkEthMainnet WebhookNetwork = "ETH_MAINNET"
	WebhookNetworkEthSepolia WebhookNetwork = "ETH_SEPOLIA"
	WebhookNetworkEthHolesky WebhookNetwork = "ETH_HOLESKY"

	// Polygon
	WebhookNetworkPolygonMainnet WebhookNetwork = "MATIC_MAINNET"
	WebhookNetworkPolygonAmoy    WebhookNetwork = "MATIC_AMOY"

	// Arbitrum
	WebhookNetworkArbitrumMainnet WebhookNetwork = "ARB_MAINNET"
	WebhookNetworkArbitrumSepolia WebhookNetwork = "ARB_SEPOLIA"

	// Optimism
	WebhookNetworkOptimismMainnet WebhookNetwork = "OPT_MAINNET"
	WebhookNetworkOptimismSepolia WebhookNetwork = "OPT_SEPOLIA"

	// Base
	WebhookNetworkBaseMainnet WebhookNetwork = "BASE_MAINNET"
	WebhookNetworkBaseSepolia WebhookNetwork = "BASE_SEPOLIA"

	// zkSync
	WebhookNetworkZkSyncMainnet WebhookNetwork = "ZKSYNC_MAINNET"
	WebhookNetworkZkSyncSepolia WebhookNetwork = "ZKSYNC_SEPOLIA"
)

// WebhookVersion represents the version of webhook payload format.
type WebhookVersion string

// Webhook versions.
const (
	WebhookVersionV1 WebhookVersion = "V1"
	WebhookVersionV2 WebhookVersion = "V2"
)

// Webhook represents an Alchemy webhook configuration.
type Webhook struct {
	// ID is the unique identifier for the webhook.
	ID string `json:"id"`
	// Network is the blockchain network.
	Network WebhookNetwork `json:"network"`
	// WebhookType is the type of webhook.
	WebhookType WebhookType `json:"webhook_type"`
	// WebhookURL is the URL where webhook events are sent.
	WebhookURL string `json:"webhook_url"`
	// IsActive indicates whether the webhook is active.
	IsActive bool `json:"is_active"`
	// TimeCreated is the Unix timestamp when the webhook was created.
	TimeCreated int64 `json:"time_created"`
	// Version is the webhook payload version (V1 or V2).
	Version WebhookVersion `json:"version"`
	// SigningKey is the key used to sign webhook payloads.
	SigningKey string `json:"signing_key"`
	// AppID is the associated app ID (optional).
	AppID *string `json:"app_id,omitempty"`
	// Name is the webhook name (optional).
	Name *string `json:"name,omitempty"`
}

// GetWebhooksResponse represents the response from the team-webhooks endpoint.
type GetWebhooksResponse struct {
	// Data contains the list of webhooks.
	Data []Webhook `json:"data"`
}

// CreateWebhookParams represents the parameters for creating a webhook.
type CreateWebhookParams struct {
	// Network is the blockchain network.
	Network WebhookNetwork `json:"network"`
	// WebhookType is the type of webhook to create.
	WebhookType WebhookType `json:"webhook_type"`
	// WebhookURL is the URL where webhook events will be sent.
	WebhookURL string `json:"webhook_url"`
	// Addresses is the list of addresses to track (for ADDRESS_ACTIVITY webhooks).
	Addresses []string `json:"addresses,omitempty"`
	// NFTFilters is the list of NFT filters (for NFT_ACTIVITY webhooks).
	NFTFilters []NFTWebhookFilter `json:"nft_filters,omitempty"`
	// GraphQLQuery is the GraphQL query (for GRAPHQL webhooks).
	GraphQLQuery *string `json:"graphql_query,omitempty"`
	// AppID is the app ID to associate with the webhook (optional).
	AppID *string `json:"app_id,omitempty"`
}

// NewAddressActivityWebhookParams creates parameters for an ADDRESS_ACTIVITY webhook.
func NewAddressActivityWebhookParams(network WebhookNetwork, webhookURL string, addresses []string) *CreateWebhookParams {
	return &CreateWebhookParams{
		Network:     network,
		WebhookType: WebhookTypeAddressActivity,
		WebhookURL:  webhookURL,
		Addresses:   addresses,
	}
}

// NewNFTActivityWebhookParams creates parameters for an NFT_ACTIVITY webhook.
func NewNFTActivityWebhookParams(network WebhookNetwork, webhookURL string, filters []NFTWebhookFilter) *CreateWebhookParams {
	return &CreateWebhookParams{
		Network:     network,
		WebhookType: WebhookTypeNFTActivity,
		WebhookURL:  webhookURL,
		NFTFilters:  filters,
	}
}

// NewGraphQLWebhookParams creates parameters for a GRAPHQL (custom) webhook.
func NewGraphQLWebhookParams(network WebhookNetwork, webhookURL string, query string) *CreateWebhookParams {
	return &CreateWebhookParams{
		Network:      network,
		WebhookType:  WebhookTypeGraphQL,
		WebhookURL:   webhookURL,
		GraphQLQuery: &query,
	}
}

// NFTWebhookFilter represents a filter for NFT activity webhooks.
type NFTWebhookFilter struct {
	// ContractAddress is the NFT contract address to track.
	ContractAddress string `json:"contract_address"`
	// TokenID is the specific token ID to track (optional, tracks all if empty).
	TokenID *string `json:"token_id,omitempty"`
}

// CreateWebhookResponse represents the response from creating a webhook.
type CreateWebhookResponse struct {
	// Data contains the created webhook.
	Data Webhook `json:"data"`
}

// UpdateWebhookParams represents the parameters for updating a webhook.
type UpdateWebhookParams struct {
	// WebhookID is the ID of the webhook to update.
	WebhookID string `json:"webhook_id"`
	// IsActive sets whether the webhook is active.
	IsActive *bool `json:"is_active,omitempty"`
	// Name sets the webhook name.
	Name *string `json:"name,omitempty"`
}

// NewUpdateWebhookParams creates parameters for updating a webhook.
func NewUpdateWebhookParams(webhookID string) *UpdateWebhookParams {
	return &UpdateWebhookParams{
		WebhookID: webhookID,
	}
}

// SetActive sets the webhook active status.
func (p *UpdateWebhookParams) SetActive(active bool) *UpdateWebhookParams {
	p.IsActive = &active
	return p
}

// SetName sets the webhook name.
func (p *UpdateWebhookParams) SetName(name string) *UpdateWebhookParams {
	p.Name = &name
	return p
}

// UpdateWebhookResponse represents the response from updating a webhook.
type UpdateWebhookResponse struct {
	// Data contains the updated webhook.
	Data Webhook `json:"data"`
}

// GetWebhookAddressesParams represents the parameters for getting webhook addresses.
type GetWebhookAddressesParams struct {
	// WebhookID is the ID of the webhook.
	WebhookID string `json:"webhook_id"`
	// Limit is the maximum number of addresses to return per page.
	Limit int `json:"limit,omitempty"`
	// After is the cursor for pagination.
	After string `json:"after,omitempty"`
	// PageKey is the page key for pagination.
	PageKey string `json:"pageKey,omitempty"`
}

// GetWebhookAddressesResponse represents the response from getting webhook addresses.
type GetWebhookAddressesResponse struct {
	// Data contains the list of addresses.
	Data []string `json:"data"`
	// Pagination contains pagination information.
	Pagination WebhookPagination `json:"pagination"`
}

// HasMore returns true if there are more addresses to fetch.
func (r *GetWebhookAddressesResponse) HasMore() bool {
	return r.Pagination.Cursors.After != ""
}

// WebhookPagination represents pagination information.
type WebhookPagination struct {
	// Cursors contains pagination cursors.
	Cursors WebhookCursors `json:"cursors"`
	// TotalCount is the total number of items.
	TotalCount int `json:"total_count"`
}

// WebhookCursors represents pagination cursors.
type WebhookCursors struct {
	// After is the cursor pointing to the end of the current result set.
	After string `json:"after"`
}

// ReplaceWebhookAddressesParams represents the parameters for replacing webhook addresses.
type ReplaceWebhookAddressesParams struct {
	// WebhookID is the ID of the webhook.
	WebhookID string `json:"webhook_id"`
	// Addresses is the new list of addresses (replaces all existing addresses).
	Addresses []string `json:"addresses"`
}

// UpdateWebhookAddressesParams represents the parameters for updating webhook addresses.
type UpdateWebhookAddressesParams struct {
	// WebhookID is the ID of the webhook.
	WebhookID string `json:"webhook_id"`
	// AddressesToAdd is the list of addresses to add.
	AddressesToAdd []string `json:"addresses_to_add"`
	// AddressesToRemove is the list of addresses to remove.
	AddressesToRemove []string `json:"addresses_to_remove"`
}

// NewUpdateWebhookAddressesParams creates parameters for updating webhook addresses.
func NewUpdateWebhookAddressesParams(webhookID string) *UpdateWebhookAddressesParams {
	return &UpdateWebhookAddressesParams{
		WebhookID:         webhookID,
		AddressesToAdd:    []string{},
		AddressesToRemove: []string{},
	}
}

// AddAddresses adds addresses to track.
func (p *UpdateWebhookAddressesParams) AddAddresses(addresses ...string) *UpdateWebhookAddressesParams {
	p.AddressesToAdd = append(p.AddressesToAdd, addresses...)
	return p
}

// RemoveAddresses removes addresses from tracking.
func (p *UpdateWebhookAddressesParams) RemoveAddresses(addresses ...string) *UpdateWebhookAddressesParams {
	p.AddressesToRemove = append(p.AddressesToRemove, addresses...)
	return p
}

// NFTWebhookFiltersResponse represents the response from getting NFT webhook filters.
type NFTWebhookFiltersResponse struct {
	// Data contains the list of NFT filters.
	Data []NFTWebhookFilter `json:"data"`
	// Pagination contains pagination information.
	Pagination WebhookPagination `json:"pagination"`
}

// UpdateNFTFiltersParams represents the parameters for updating NFT webhook filters.
type UpdateNFTFiltersParams struct {
	// WebhookID is the ID of the webhook.
	WebhookID string `json:"webhook_id"`
	// FiltersToAdd is the list of filters to add.
	FiltersToAdd []NFTWebhookFilter `json:"nft_filters_to_add"`
	// FiltersToRemove is the list of filters to remove.
	FiltersToRemove []NFTWebhookFilter `json:"nft_filters_to_remove"`
}

// WebhookSignatureHeader is the HTTP header containing the webhook signature.
const WebhookSignatureHeader = "X-Alchemy-Signature"

// WebhookEvent represents a webhook event payload (V2 format).
type WebhookEvent struct {
	// WebhookID is the unique identifier for the webhook.
	WebhookID string `json:"webhookId"`
	// ID is the event ID.
	ID string `json:"id"`
	// CreatedAt is when the event was created.
	CreatedAt string `json:"createdAt"`
	// Type is the event type.
	Type string `json:"type"`
	// Event contains the event data.
	Event interface{} `json:"event"`
}

// AddressActivityEvent represents an address activity event.
type AddressActivityEvent struct {
	// Network is the blockchain network.
	Network string `json:"network"`
	// Activity contains the activity details.
	Activity []AddressActivity `json:"activity"`
}

// AddressActivity represents a single address activity.
type AddressActivity struct {
	// FromAddress is the sender address.
	FromAddress string `json:"fromAddress"`
	// ToAddress is the recipient address.
	ToAddress string `json:"toAddress"`
	// BlockNum is the block number (hex).
	BlockNum string `json:"blockNum"`
	// Hash is the transaction hash.
	Hash string `json:"hash"`
	// Value is the value transferred.
	Value float64 `json:"value"`
	// Asset is the asset symbol.
	Asset string `json:"asset"`
	// Category is the transfer category.
	Category string `json:"category"`
	// RawContract contains raw contract info.
	RawContract *RawContractInfo `json:"rawContract,omitempty"`
	// Log contains log info (for token transfers).
	Log *ActivityLog `json:"log,omitempty"`
}

// RawContractInfo contains raw contract information.
type RawContractInfo struct {
	// RawValue is the raw value (hex).
	RawValue string `json:"rawValue"`
	// Address is the contract address.
	Address string `json:"address"`
	// Decimals is the token decimals.
	Decimals int `json:"decimals"`
}

// ActivityLog contains log information for token transfers.
type ActivityLog struct {
	// Address is the contract address.
	Address string `json:"address"`
	// Topics is the list of log topics.
	Topics []string `json:"topics"`
	// Data is the log data.
	Data string `json:"data"`
	// BlockNumber is the block number (hex).
	BlockNumber string `json:"blockNumber"`
	// TransactionHash is the transaction hash.
	TransactionHash string `json:"transactionHash"`
	// TransactionIndex is the transaction index (hex).
	TransactionIndex string `json:"transactionIndex"`
	// BlockHash is the block hash.
	BlockHash string `json:"blockHash"`
	// LogIndex is the log index (hex).
	LogIndex string `json:"logIndex"`
	// Removed indicates if the log was removed.
	Removed bool `json:"removed"`
}
