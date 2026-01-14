package alchemy

// Network represents a blockchain network supported by Alchemy.
type Network string

// Ethereum networks.
const (
	EthMainnet Network = "eth-mainnet"
	EthSepolia Network = "eth-sepolia"
	EthHolesky Network = "eth-holesky"
	EthHoodi   Network = "eth-hoodi"
)

// Polygon networks.
const (
	PolygonMainnet Network = "polygon-mainnet"
	PolygonAmoy    Network = "polygon-amoy"
)

// Arbitrum networks.
const (
	ArbitrumMainnet Network = "arb-mainnet"
	ArbitrumSepolia Network = "arb-sepolia"
	ArbitrumNova    Network = "arbnova-mainnet"
)

// Optimism networks.
const (
	OptimismMainnet Network = "opt-mainnet"
	OptimismSepolia Network = "opt-sepolia"
)

// Base networks.
const (
	BaseMainnet Network = "base-mainnet"
	BaseSepolia Network = "base-sepolia"
)

// zkSync networks.
const (
	ZkSyncMainnet Network = "zksync-mainnet"
	ZkSyncSepolia Network = "zksync-sepolia"
)

// Polygon zkEVM networks.
const (
	PolygonZkEvmMainnet Network = "polygonzkevm-mainnet"
	PolygonZkEvmCardona Network = "polygonzkevm-cardona"
)

// Linea networks.
const (
	LineaMainnet Network = "linea-mainnet"
	LineaSepolia Network = "linea-sepolia"
)

// Scroll networks.
const (
	ScrollMainnet Network = "scroll-mainnet"
	ScrollSepolia Network = "scroll-sepolia"
)

// Blast networks.
const (
	BlastMainnet Network = "blast-mainnet"
	BlastSepolia Network = "blast-sepolia"
)

// Other networks.
const (
	AvalancheMainnet  Network = "avax-mainnet"
	AvalancheFuji     Network = "avax-fuji"
	BNBMainnet        Network = "bnb-mainnet"
	BNBTestnet        Network = "bnb-testnet"
	FantomMainnet     Network = "fantom-mainnet"
	FantomTestnet     Network = "fantom-testnet"
	GnosisMainnet     Network = "gnosis-mainnet"
	GnosisChiado      Network = "gnosis-chiado"
	CeloMainnet       Network = "celo-mainnet"
	CeloAlfajores     Network = "celo-alfajores"
	MantleMainnet     Network = "mantle-mainnet"
	MantleSepolia     Network = "mantle-sepolia"
	WorldChainMainnet Network = "worldchain-mainnet"
	WorldChainSepolia Network = "worldchain-sepolia"
	ZoraMainnet       Network = "zora-mainnet"
	ZoraSepolia       Network = "zora-sepolia"
	BerachainBartio   Network = "berachain-bartio"
	FlowMainnet       Network = "flow-mainnet"
	FlowTestnet       Network = "flow-testnet"
)

// BaseURL returns the base URL for the network's API endpoint.
func (n Network) BaseURL() string {
	return "https://" + string(n) + ".g.alchemy.com/v2"
}

// NFTURL returns the NFT API base URL for the network.
func (n Network) NFTURL() string {
	return "https://" + string(n) + ".g.alchemy.com/nft/v3"
}

// String returns the network identifier string.
func (n Network) String() string {
	return string(n)
}

// ChainID returns the chain ID for the network.
// Returns 0 if the network is unknown.
func (n Network) ChainID() uint64 {
	switch n {
	// Ethereum
	case EthMainnet:
		return 1
	case EthSepolia:
		return 11155111
	case EthHolesky:
		return 17000
	case EthHoodi:
		return 560048

	// Polygon
	case PolygonMainnet:
		return 137
	case PolygonAmoy:
		return 80002

	// Arbitrum
	case ArbitrumMainnet:
		return 42161
	case ArbitrumSepolia:
		return 421614
	case ArbitrumNova:
		return 42170

	// Optimism
	case OptimismMainnet:
		return 10
	case OptimismSepolia:
		return 11155420

	// Base
	case BaseMainnet:
		return 8453
	case BaseSepolia:
		return 84532

	// zkSync
	case ZkSyncMainnet:
		return 324
	case ZkSyncSepolia:
		return 300

	// Polygon zkEVM
	case PolygonZkEvmMainnet:
		return 1101
	case PolygonZkEvmCardona:
		return 2442

	// Linea
	case LineaMainnet:
		return 59144
	case LineaSepolia:
		return 59141

	// Scroll
	case ScrollMainnet:
		return 534352
	case ScrollSepolia:
		return 534351

	// Blast
	case BlastMainnet:
		return 81457
	case BlastSepolia:
		return 168587773

	// Avalanche
	case AvalancheMainnet:
		return 43114
	case AvalancheFuji:
		return 43113

	// BNB
	case BNBMainnet:
		return 56
	case BNBTestnet:
		return 97

	// Fantom
	case FantomMainnet:
		return 250
	case FantomTestnet:
		return 4002

	// Gnosis
	case GnosisMainnet:
		return 100
	case GnosisChiado:
		return 10200

	// Celo
	case CeloMainnet:
		return 42220
	case CeloAlfajores:
		return 44787

	// Mantle
	case MantleMainnet:
		return 5000
	case MantleSepolia:
		return 5003

	// World Chain
	case WorldChainMainnet:
		return 480
	case WorldChainSepolia:
		return 4801

	// Zora
	case ZoraMainnet:
		return 7777777
	case ZoraSepolia:
		return 999999999

	// Berachain
	case BerachainBartio:
		return 80084

	// Flow
	case FlowMainnet:
		return 747
	case FlowTestnet:
		return 545

	default:
		return 0
	}
}

// IsMainnet returns true if this is a mainnet network.
func (n Network) IsMainnet() bool {
	switch n {
	case EthMainnet, PolygonMainnet, ArbitrumMainnet, ArbitrumNova,
		OptimismMainnet, BaseMainnet, ZkSyncMainnet, PolygonZkEvmMainnet,
		LineaMainnet, ScrollMainnet, BlastMainnet, AvalancheMainnet,
		BNBMainnet, FantomMainnet, GnosisMainnet, CeloMainnet,
		MantleMainnet, WorldChainMainnet, ZoraMainnet, FlowMainnet:
		return true
	default:
		return false
	}
}

// IsTestnet returns true if this is a testnet network.
func (n Network) IsTestnet() bool {
	return !n.IsMainnet()
}

// IsEthereum returns true if this is an Ethereum network.
func (n Network) IsEthereum() bool {
	switch n {
	case EthMainnet, EthSepolia, EthHolesky, EthHoodi:
		return true
	default:
		return false
	}
}

// IsL2 returns true if this is a Layer 2 network.
func (n Network) IsL2() bool {
	switch n {
	case ArbitrumMainnet, ArbitrumSepolia, ArbitrumNova,
		OptimismMainnet, OptimismSepolia,
		BaseMainnet, BaseSepolia,
		ZkSyncMainnet, ZkSyncSepolia,
		PolygonZkEvmMainnet, PolygonZkEvmCardona,
		LineaMainnet, LineaSepolia,
		ScrollMainnet, ScrollSepolia,
		BlastMainnet, BlastSepolia,
		MantleMainnet, MantleSepolia,
		ZoraMainnet, ZoraSepolia:
		return true
	default:
		return false
	}
}

// NativeCurrency returns the native currency symbol for the network.
func (n Network) NativeCurrency() string {
	switch n {
	case PolygonMainnet, PolygonAmoy:
		return "MATIC"
	case AvalancheMainnet, AvalancheFuji:
		return "AVAX"
	case BNBMainnet, BNBTestnet:
		return "BNB"
	case FantomMainnet, FantomTestnet:
		return "FTM"
	case GnosisMainnet, GnosisChiado:
		return "xDAI"
	case CeloMainnet, CeloAlfajores:
		return "CELO"
	case MantleMainnet, MantleSepolia:
		return "MNT"
	case FlowMainnet, FlowTestnet:
		return "FLOW"
	case BerachainBartio:
		return "BERA"
	default:
		return "ETH"
	}
}

// AllNetworks returns a list of all supported networks.
func AllNetworks() []Network {
	return []Network{
		// Ethereum
		EthMainnet, EthSepolia, EthHolesky, EthHoodi,
		// Polygon
		PolygonMainnet, PolygonAmoy,
		// Arbitrum
		ArbitrumMainnet, ArbitrumSepolia, ArbitrumNova,
		// Optimism
		OptimismMainnet, OptimismSepolia,
		// Base
		BaseMainnet, BaseSepolia,
		// zkSync
		ZkSyncMainnet, ZkSyncSepolia,
		// Polygon zkEVM
		PolygonZkEvmMainnet, PolygonZkEvmCardona,
		// Linea
		LineaMainnet, LineaSepolia,
		// Scroll
		ScrollMainnet, ScrollSepolia,
		// Blast
		BlastMainnet, BlastSepolia,
		// Others
		AvalancheMainnet, AvalancheFuji,
		BNBMainnet, BNBTestnet,
		FantomMainnet, FantomTestnet,
		GnosisMainnet, GnosisChiado,
		CeloMainnet, CeloAlfajores,
		MantleMainnet, MantleSepolia,
		WorldChainMainnet, WorldChainSepolia,
		ZoraMainnet, ZoraSepolia,
		BerachainBartio,
		FlowMainnet, FlowTestnet,
	}
}

// MainnetNetworks returns a list of all mainnet networks.
func MainnetNetworks() []Network {
	var networks []Network
	for _, n := range AllNetworks() {
		if n.IsMainnet() {
			networks = append(networks, n)
		}
	}
	return networks
}
