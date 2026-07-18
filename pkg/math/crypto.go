package math

import (
	"crypto/sha256"
	"encoding/binary"
)

func CryptoNumberFromStringSHA(input string, max int) int {
	hash := sha256.Sum256([]byte(input))

	num := binary.BigEndian.Uint64(hash[:8])

	return int(num%uint64(max)) + 1
}
