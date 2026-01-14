// Package hex provides utilities for encoding and decoding hexadecimal values
// commonly used in Ethereum JSON-RPC APIs.
package hex

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// Has0xPrefix returns true if s has a 0x prefix.
func Has0xPrefix(s string) bool {
	return len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')
}

// TrimPrefix removes the 0x prefix from s if present.
func TrimPrefix(s string) string {
	if Has0xPrefix(s) {
		return s[2:]
	}
	return s
}

// AddPrefix adds 0x prefix to s if not present.
func AddPrefix(s string) string {
	if Has0xPrefix(s) {
		return s
	}
	return "0x" + s
}

// Encode encodes bytes as a hex string with 0x prefix.
func Encode(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

// Decode decodes a hex string (with or without 0x prefix) to bytes.
func Decode(s string) ([]byte, error) {
	s = TrimPrefix(s)
	if len(s)%2 != 0 {
		s = "0" + s
	}
	return hex.DecodeString(s)
}

// MustDecode is like Decode but panics on error.
func MustDecode(s string) []byte {
	b, err := Decode(s)
	if err != nil {
		panic(err)
	}
	return b
}

// EncodeUint64 encodes a uint64 as a hex string with 0x prefix.
// The encoding uses the minimum number of hex digits (no leading zeros).
func EncodeUint64(n uint64) string {
	if n == 0 {
		return "0x0"
	}
	return "0x" + strconv.FormatUint(n, 16)
}

// DecodeUint64 decodes a hex string to uint64.
func DecodeUint64(s string) (uint64, error) {
	s = TrimPrefix(s)
	if s == "" {
		return 0, nil
	}
	return strconv.ParseUint(s, 16, 64)
}

// MustDecodeUint64 is like DecodeUint64 but panics on error.
func MustDecodeUint64(s string) uint64 {
	n, err := DecodeUint64(s)
	if err != nil {
		panic(err)
	}
	return n
}

// EncodeBigInt encodes a big.Int as a hex string with 0x prefix.
func EncodeBigInt(n *big.Int) string {
	if n == nil || n.Sign() == 0 {
		return "0x0"
	}
	if n.Sign() < 0 {
		// Negative numbers are encoded as two's complement
		return "-0x" + n.Text(16)[1:]
	}
	return "0x" + n.Text(16)
}

// DecodeBigInt decodes a hex string to *big.Int.
func DecodeBigInt(s string) (*big.Int, error) {
	if s == "" || s == "0x" || s == "0x0" {
		return big.NewInt(0), nil
	}

	negative := false
	if strings.HasPrefix(s, "-") {
		negative = true
		s = s[1:]
	}

	s = TrimPrefix(s)
	n, ok := new(big.Int).SetString(s, 16)
	if !ok {
		return nil, fmt.Errorf("invalid hex big int: %s", s)
	}

	if negative {
		n.Neg(n)
	}
	return n, nil
}

// MustDecodeBigInt is like DecodeBigInt but panics on error.
func MustDecodeBigInt(s string) *big.Int {
	n, err := DecodeBigInt(s)
	if err != nil {
		panic(err)
	}
	return n
}

// IsValidHex checks if s is a valid hex string (with 0x prefix).
func IsValidHex(s string) bool {
	if !Has0xPrefix(s) {
		return false
	}
	s = s[2:]
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if !isHexDigit(byte(c)) {
			return false
		}
	}
	return true
}

// IsValidAddress checks if s is a valid Ethereum address (40 hex chars with 0x prefix).
func IsValidAddress(s string) bool {
	if !Has0xPrefix(s) {
		return false
	}
	s = s[2:]
	if len(s) != 40 {
		return false
	}
	for _, c := range s {
		if !isHexDigit(byte(c)) {
			return false
		}
	}
	return true
}

// IsValidHash checks if s is a valid Ethereum hash (64 hex chars with 0x prefix).
func IsValidHash(s string) bool {
	if !Has0xPrefix(s) {
		return false
	}
	s = s[2:]
	if len(s) != 64 {
		return false
	}
	for _, c := range s {
		if !isHexDigit(byte(c)) {
			return false
		}
	}
	return true
}

func isHexDigit(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
