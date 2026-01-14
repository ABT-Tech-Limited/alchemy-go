package types

// Log represents an Ethereum log entry.
type Log struct {
	// Address is the contract address that emitted the log.
	Address Address `json:"address"`
	// Topics is the list of topics (indexed parameters).
	Topics []Hash `json:"topics"`
	// Data is the non-indexed log data.
	Data Data `json:"data"`
	// BlockNumber is the block number where the log was emitted.
	BlockNumber Quantity `json:"blockNumber"`
	// TransactionHash is the transaction hash.
	TransactionHash Hash `json:"transactionHash"`
	// TransactionIndex is the transaction index in the block.
	TransactionIndex Quantity `json:"transactionIndex"`
	// BlockHash is the block hash.
	BlockHash Hash `json:"blockHash"`
	// LogIndex is the log index in the block.
	LogIndex Quantity `json:"logIndex"`
	// Removed is true if the log was reverted due to a chain reorganization.
	Removed bool `json:"removed"`
}

// Topic0 returns the first topic (event signature), or empty hash if not present.
func (l *Log) Topic0() Hash {
	if len(l.Topics) == 0 {
		return ""
	}
	return l.Topics[0]
}

// Topic1 returns the second topic, or empty hash if not present.
func (l *Log) Topic1() Hash {
	if len(l.Topics) < 2 {
		return ""
	}
	return l.Topics[1]
}

// Topic2 returns the third topic, or empty hash if not present.
func (l *Log) Topic2() Hash {
	if len(l.Topics) < 3 {
		return ""
	}
	return l.Topics[2]
}

// Topic3 returns the fourth topic, or empty hash if not present.
func (l *Log) Topic3() Hash {
	if len(l.Topics) < 4 {
		return ""
	}
	return l.Topics[3]
}
