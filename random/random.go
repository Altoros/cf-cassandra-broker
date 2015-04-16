package random

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

func Bytes(n uint) []byte {
	bytes := make([]byte, n)

	for i := 0; i < int(n); i++ {
		randomByte, err := rand.Int(rand.Reader, big.NewInt(255))
		if err != nil {
			panic(err.Error())
		}
		bytes[i] = byte(randomByte.Int64())
	}

	return bytes
}

func Hex(n uint) string {
	return hex.EncodeToString(Bytes(n))
}
