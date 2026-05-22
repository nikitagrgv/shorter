package token_hasher

import (
	"crypto/sha256"
	"encoding/hex"
)

type TokenHasher struct{}

func (t TokenHasher) Hash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}
