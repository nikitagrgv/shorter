package domain

type TokenGenerator interface {
	Generate() string
}
