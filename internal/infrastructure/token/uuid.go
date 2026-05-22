package token

import (
	"github.com/google/uuid"
)

type UuidTokenGenerator struct{}

func (UuidTokenGenerator) Generate() string {
	return uuid.NewString()
}
