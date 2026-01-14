// Package types provides common types used across the Alchemy SDK.
package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ABT-Tech-Limited/alchemy-go/internal/hex"
)

// Address represents an Ethereum address (20 bytes).
type Address string

// ZeroAddress is the zero address.
const ZeroAddress Address = "0x0000000000000000000000000000000000000000"

// ParseAddress parses and validates an Ethereum address string.
func ParseAddress(s string) (Address, error) {
	s = strings.ToLower(s)
	if !hex.Has0xPrefix(s) {
		s = "0x" + s
	}
	if !hex.IsValidAddress(s) {
		return "", fmt.Errorf("invalid address: %s", s)
	}
	return Address(s), nil
}

// MustParseAddress is like ParseAddress but panics on error.
func MustParseAddress(s string) Address {
	addr, err := ParseAddress(s)
	if err != nil {
		panic(err)
	}
	return addr
}

// String returns the hex string representation of the address.
func (a Address) String() string {
	return string(a)
}

// Bytes returns the address as a 20-byte slice.
func (a Address) Bytes() []byte {
	if a == "" {
		return nil
	}
	b, _ := hex.Decode(string(a))
	return b
}

// IsZero returns true if the address is the zero address or empty.
func (a Address) IsZero() bool {
	return a == "" || a == ZeroAddress
}

// MarshalJSON implements json.Marshaler.
func (a Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(a))
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *Address) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		*a = ""
		return nil
	}
	addr, err := ParseAddress(s)
	if err != nil {
		return err
	}
	*a = addr
	return nil
}

// Hash represents an Ethereum hash (32 bytes).
type Hash string

// ZeroHash is the zero hash.
const ZeroHash Hash = "0x0000000000000000000000000000000000000000000000000000000000000000"

// ParseHash parses and validates an Ethereum hash string.
func ParseHash(s string) (Hash, error) {
	s = strings.ToLower(s)
	if !hex.Has0xPrefix(s) {
		s = "0x" + s
	}
	if !hex.IsValidHash(s) {
		return "", fmt.Errorf("invalid hash: %s", s)
	}
	return Hash(s), nil
}

// MustParseHash is like ParseHash but panics on error.
func MustParseHash(s string) Hash {
	h, err := ParseHash(s)
	if err != nil {
		panic(err)
	}
	return h
}

// String returns the hex string representation of the hash.
func (h Hash) String() string {
	return string(h)
}

// Bytes returns the hash as a 32-byte slice.
func (h Hash) Bytes() []byte {
	if h == "" {
		return nil
	}
	b, _ := hex.Decode(string(h))
	return b
}

// IsZero returns true if the hash is the zero hash or empty.
func (h Hash) IsZero() bool {
	return h == "" || h == ZeroHash
}

// MarshalJSON implements json.Marshaler.
func (h Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(h))
}

// UnmarshalJSON implements json.Unmarshaler.
func (h *Hash) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		*h = ""
		return nil
	}
	hash, err := ParseHash(s)
	if err != nil {
		// Some APIs may return non-standard hashes, accept them
		*h = Hash(strings.ToLower(s))
		return nil
	}
	*h = hash
	return nil
}

// Quantity represents a hex-encoded quantity (used for numbers in JSON-RPC).
type Quantity string

// QuantityFromUint64 creates a Quantity from a uint64.
func QuantityFromUint64(n uint64) Quantity {
	return Quantity(hex.EncodeUint64(n))
}

// QuantityFromBigInt creates a Quantity from a *big.Int.
func QuantityFromBigInt(n *big.Int) Quantity {
	if n == nil {
		return "0x0"
	}
	return Quantity(hex.EncodeBigInt(n))
}

// Uint64 returns the quantity as uint64.
func (q Quantity) Uint64() uint64 {
	n, _ := hex.DecodeUint64(string(q))
	return n
}

// BigInt returns the quantity as *big.Int.
func (q Quantity) BigInt() *big.Int {
	n, _ := hex.DecodeBigInt(string(q))
	if n == nil {
		return big.NewInt(0)
	}
	return n
}

// Int64 returns the quantity as int64.
func (q Quantity) Int64() int64 {
	return q.BigInt().Int64()
}

// String returns the hex string representation.
func (q Quantity) String() string {
	return string(q)
}

// IsZero returns true if the quantity is zero or empty.
func (q Quantity) IsZero() bool {
	return q == "" || q == "0x0" || q == "0x"
}

// MarshalJSON implements json.Marshaler.
func (q Quantity) MarshalJSON() ([]byte, error) {
	if q == "" {
		return json.Marshal("0x0")
	}
	return json.Marshal(string(q))
}

// UnmarshalJSON implements json.Unmarshaler.
func (q *Quantity) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		// Try to unmarshal as number
		var n uint64
		if err := json.Unmarshal(data, &n); err != nil {
			return err
		}
		*q = QuantityFromUint64(n)
		return nil
	}
	*q = Quantity(s)
	return nil
}

// Data represents arbitrary hex-encoded data.
type Data string

// DataFromBytes creates Data from bytes.
func DataFromBytes(b []byte) Data {
	return Data(hex.Encode(b))
}

// Bytes returns the data as bytes.
func (d Data) Bytes() []byte {
	b, _ := hex.Decode(string(d))
	return b
}

// String returns the hex string representation.
func (d Data) String() string {
	return string(d)
}

// MarshalJSON implements json.Marshaler.
func (d Data) MarshalJSON() ([]byte, error) {
	if d == "" {
		return json.Marshal("0x")
	}
	return json.Marshal(string(d))
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *Data) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*d = Data(s)
	return nil
}
