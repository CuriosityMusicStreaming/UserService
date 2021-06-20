//nolint:gosec
package hash

import (
	"crypto/sha1"
	"fmt"
)

type Hasher interface {
	Hash(value string) string
}

type sha1Hasher struct {
	salt string
}

func NewSHA1Hasher(salt string) Hasher {
	return &sha1Hasher{salt: salt}
}

func (h *sha1Hasher) Hash(value string) string {
	hash := sha1.New()
	_, _ = hash.Write([]byte(value))

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt)))
}
