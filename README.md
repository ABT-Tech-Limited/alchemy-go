# Alchemy Go SDK

Unofficial Go SDK for [Alchemy](https://www.alchemy.com/) blockchain API.

## Installation

```bash
go get github.com/ABT-Tech-Limited/alchemy-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/ABT-Tech-Limited/alchemy-go"
    "github.com/ABT-Tech-Limited/alchemy-go/data"
    "github.com/ABT-Tech-Limited/alchemy-go/types"
)

func main() {
    client, err := alchemy.New(alchemy.Config{
        APIKey:  os.Getenv("ALCHEMY_API_KEY"),
        Network: alchemy.EthMainnet,
    })
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Get current block number
    blockNum, _ := client.Node.BlockNumber(ctx)
    fmt.Printf("Block: %d\n", blockNum)

    // Get wallet balance
    addr := types.MustParseAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
    balance, _ := client.Wallet.GetBalance(ctx, addr)
    fmt.Printf("Balance: %s ETH\n", balance.Formatted)

    // Get asset transfers
    params := data.NewAssetTransfersParams().
        SetToAddress(addr).
        SetCategories([]data.AssetTransferCategory{data.CategoryERC20})
    transfers, _ := client.Data.GetAssetTransfers(ctx, params)
    fmt.Printf("Transfers: %d\n", len(transfers.Transfers))
}
```

## Supported Networks

| Network | Constant |
|---------|----------|
| Ethereum Mainnet | `alchemy.EthMainnet` |
| Ethereum Sepolia | `alchemy.EthSepolia` |
| Polygon Mainnet | `alchemy.PolygonMainnet` |
| Arbitrum Mainnet | `alchemy.ArbitrumMainnet` |
| Optimism Mainnet | `alchemy.OptimismMainnet` |
| Base Mainnet | `alchemy.BaseMainnet` |

## License

MIT
