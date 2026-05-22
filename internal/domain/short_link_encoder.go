package domain

type ShortLinkEncoder interface {
	Encode(id int64) string
	Decode(code string) (int64, error)
}
