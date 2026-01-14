package data

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// WebhookClient provides access to Alchemy Webhook (Notify) API.
// It requires an auth token obtained from the Alchemy dashboard.
type WebhookClient struct {
	authToken  string
	httpClient *http.Client
	baseURL    string
}

// NewWebhookClient creates a new WebhookClient.
func NewWebhookClient(authToken string, httpClient *http.Client) *WebhookClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &WebhookClient{
		authToken:  authToken,
		httpClient: httpClient,
		baseURL:    "https://dashboard.alchemy.com/api",
	}
}

// GetAllWebhooks retrieves all webhooks for the team.
func (c *WebhookClient) GetAllWebhooks(ctx context.Context) (*GetWebhooksResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/team-webhooks", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}

	var result GetWebhooksResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// CreateWebhook creates a new webhook.
func (c *WebhookClient) CreateWebhook(ctx context.Context, params *CreateWebhookParams) (*CreateWebhookResponse, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/create-webhook", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}

	var result CreateWebhookResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// UpdateWebhook updates a webhook's status or name.
func (c *WebhookClient) UpdateWebhook(ctx context.Context, params *UpdateWebhookParams) (*UpdateWebhookResponse, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.baseURL+"/update-webhook", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}

	var result UpdateWebhookResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// DeleteWebhook deletes a webhook.
func (c *WebhookClient) DeleteWebhook(ctx context.Context, webhookID string) error {
	reqURL := c.baseURL + "/delete-webhook?webhook_id=" + url.QueryEscape(webhookID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.setAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	return c.checkResponse(resp)
}

// GetWebhookAddresses retrieves addresses tracked by a webhook.
func (c *WebhookClient) GetWebhookAddresses(ctx context.Context, params *GetWebhookAddressesParams) (*GetWebhookAddressesResponse, error) {
	query := url.Values{}
	query.Set("webhook_id", params.WebhookID)

	if params.Limit > 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.After != "" {
		query.Set("after", params.After)
	}
	if params.PageKey != "" {
		query.Set("pageKey", params.PageKey)
	}

	reqURL := c.baseURL + "/webhook-addresses?" + query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}

	var result GetWebhookAddressesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetAllWebhookAddresses retrieves all addresses tracked by a webhook (handles pagination).
func (c *WebhookClient) GetAllWebhookAddresses(ctx context.Context, webhookID string) ([]string, error) {
	var allAddresses []string
	after := ""

	for {
		params := &GetWebhookAddressesParams{
			WebhookID: webhookID,
			Limit:     1000,
			After:     after,
		}

		resp, err := c.GetWebhookAddresses(ctx, params)
		if err != nil {
			return nil, err
		}

		allAddresses = append(allAddresses, resp.Data...)

		if !resp.HasMore() {
			break
		}
		after = resp.Pagination.Cursors.After
	}

	return allAddresses, nil
}

// ReplaceWebhookAddresses replaces all addresses tracked by a webhook.
func (c *WebhookClient) ReplaceWebhookAddresses(ctx context.Context, params *ReplaceWebhookAddressesParams) error {
	body, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.baseURL+"/update-webhook-addresses", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.setAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	return c.checkResponse(resp)
}

// UpdateWebhookAddresses adds or removes addresses from a webhook.
func (c *WebhookClient) UpdateWebhookAddresses(ctx context.Context, params *UpdateWebhookAddressesParams) error {
	body, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.baseURL+"/update-webhook-addresses", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.setAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	return c.checkResponse(resp)
}

// GetNFTFilters retrieves NFT filters for a webhook.
func (c *WebhookClient) GetNFTFilters(ctx context.Context, webhookID string, limit int, after string) (*NFTWebhookFiltersResponse, error) {
	query := url.Values{}
	query.Set("webhook_id", webhookID)

	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}
	if after != "" {
		query.Set("after", after)
	}

	reqURL := c.baseURL + "/webhook-nft-filters?" + query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}

	var result NFTWebhookFiltersResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// UpdateNFTFilters adds or removes NFT filters from a webhook.
func (c *WebhookClient) UpdateNFTFilters(ctx context.Context, params *UpdateNFTFiltersParams) error {
	body, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.baseURL+"/update-webhook-nft-filters", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.setAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	return c.checkResponse(resp)
}

// setAuthHeader sets the authentication header.
func (c *WebhookClient) setAuthHeader(req *http.Request) {
	req.Header.Set("X-Alchemy-Token", c.authToken)
}

// checkResponse checks the HTTP response for errors.
func (c *WebhookClient) checkResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(body))
}

// VerifyWebhookSignature verifies the signature of a webhook payload.
// signingKey is the webhook's signing key from the dashboard.
// signature is the value of the X-Alchemy-Signature header.
// payload is the raw request body.
func VerifyWebhookSignature(signingKey, signature string, payload []byte) bool {
	mac := hmac.New(sha256.New, []byte(signingKey))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// ParseWebhookEvent parses a webhook event from the request body.
func ParseWebhookEvent(body []byte) (*WebhookEvent, error) {
	var event WebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, fmt.Errorf("failed to parse webhook event: %w", err)
	}
	return &event, nil
}

// ParseAddressActivityEvent parses the event data as an AddressActivityEvent.
func ParseAddressActivityEvent(event *WebhookEvent) (*AddressActivityEvent, error) {
	data, err := json.Marshal(event.Event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event data: %w", err)
	}

	var activity AddressActivityEvent
	if err := json.Unmarshal(data, &activity); err != nil {
		return nil, fmt.Errorf("failed to parse address activity event: %w", err)
	}

	return &activity, nil
}
