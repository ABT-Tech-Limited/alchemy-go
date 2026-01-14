package data

import (
	"context"
	"sync"
)

// GetAssetTransfers retrieves asset transfers matching the given parameters.
func (c *Client) GetAssetTransfers(ctx context.Context, params *AssetTransfersParams) (*AssetTransfersResponse, error) {
	var result AssetTransfersResponse
	if err := c.rpc.Call(ctx, "alchemy_getAssetTransfers", []interface{}{params}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAssetTransfersIterator returns an iterator for paginating through asset transfers.
func (c *Client) GetAssetTransfersIterator(ctx context.Context, params *AssetTransfersParams) *AssetTransfersIterator {
	// Make a copy of params to avoid modifying the original
	paramsCopy := *params
	return &AssetTransfersIterator{
		client: c,
		params: &paramsCopy,
		ctx:    ctx,
	}
}

// AssetTransfersIterator iterates through asset transfers with pagination.
type AssetTransfersIterator struct {
	client  *Client
	params  *AssetTransfersParams
	ctx     context.Context
	current *AssetTransfersResponse
	index   int
	done    bool
	err     error
	mu      sync.Mutex
}

// Next returns the next transfer in the iteration.
// Returns nil when there are no more transfers.
func (it *AssetTransfersIterator) Next() (*AssetTransfer, error) {
	it.mu.Lock()
	defer it.mu.Unlock()

	if it.err != nil {
		return nil, it.err
	}

	if it.done {
		return nil, nil
	}

	// Fetch first page if needed
	if it.current == nil {
		if err := it.fetchNext(); err != nil {
			it.err = err
			return nil, err
		}
	}

	// Check if we have more transfers in current page
	if it.index < len(it.current.Transfers) {
		transfer := &it.current.Transfers[it.index]
		it.index++
		return transfer, nil
	}

	// Check if there are more pages
	if !it.current.HasMore() {
		it.done = true
		return nil, nil
	}

	// Fetch next page
	it.params.PageKey = it.current.PageKey
	if err := it.fetchNext(); err != nil {
		it.err = err
		return nil, err
	}

	if len(it.current.Transfers) == 0 {
		it.done = true
		return nil, nil
	}

	transfer := &it.current.Transfers[0]
	it.index = 1
	return transfer, nil
}

// HasNext returns true if there are more transfers to iterate.
func (it *AssetTransfersIterator) HasNext() bool {
	it.mu.Lock()
	defer it.mu.Unlock()

	if it.done || it.err != nil {
		return false
	}

	// Check if we have more in current page
	if it.current != nil && it.index < len(it.current.Transfers) {
		return true
	}

	// Check if there are more pages
	if it.current != nil {
		return it.current.HasMore()
	}

	// Haven't fetched first page yet
	return true
}

// Error returns any error encountered during iteration.
func (it *AssetTransfersIterator) Error() error {
	it.mu.Lock()
	defer it.mu.Unlock()
	return it.err
}

// Reset resets the iterator to the beginning.
func (it *AssetTransfersIterator) Reset() {
	it.mu.Lock()
	defer it.mu.Unlock()

	it.current = nil
	it.index = 0
	it.done = false
	it.err = nil
	it.params.PageKey = ""
}

// Collect returns all remaining transfers as a slice.
// Use with caution on large result sets.
func (it *AssetTransfersIterator) Collect() ([]AssetTransfer, error) {
	var transfers []AssetTransfer

	for {
		transfer, err := it.Next()
		if err != nil {
			return nil, err
		}
		if transfer == nil {
			break
		}
		transfers = append(transfers, *transfer)
	}

	return transfers, nil
}

// CollectN returns up to n transfers.
func (it *AssetTransfersIterator) CollectN(n int) ([]AssetTransfer, error) {
	transfers := make([]AssetTransfer, 0, n)

	for i := 0; i < n; i++ {
		transfer, err := it.Next()
		if err != nil {
			return nil, err
		}
		if transfer == nil {
			break
		}
		transfers = append(transfers, *transfer)
	}

	return transfers, nil
}

func (it *AssetTransfersIterator) fetchNext() error {
	result, err := it.client.GetAssetTransfers(it.ctx, it.params)
	if err != nil {
		return err
	}
	it.current = result
	it.index = 0
	return nil
}
