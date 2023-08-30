package selector

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func Select(values []uint64, count int) ([]uint64, error) {
	n := len(values)
	if count <= 0 || n < count {
		return nil, fmt.Errorf("unexpected count")
	}
	for i := n - 1; i >= n-count; i-- {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return nil, err
		}
		j := int(idx.Int64())
		values[j], values[i] = values[i], values[j]
	}
	return values[(n - count):], nil
}
