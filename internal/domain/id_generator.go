package domain

type IdGenerator interface {
	Generate() int64
}
