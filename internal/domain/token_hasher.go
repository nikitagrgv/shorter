package domain

type TokenHasher interface {
	Hash(string) string
}
