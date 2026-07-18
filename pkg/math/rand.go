package math

import (
	"hash/fnv"
	"math/rand"
)

func RandNumberFromString(input string, max int) int {
	h := fnv.New64a()
	h.Write([]byte(input))
	seed := h.Sum64()

	r := rand.New(rand.NewSource(int64(seed)))

	return r.Intn(max) + 1
}
