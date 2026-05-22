package shortlink

import (
	"errors"

	"github.com/speps/go-hashids/v2"
)

type HashidsLinkEncoder struct {
	h *hashids.HashID
}

func NewHashidsLinkEncoder() *HashidsLinkEncoder {
	hd := hashids.NewData()
	hd.Salt = "3a1b3e9d"
	hd.MinLength = 0
	hd.Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-_"
	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}
	return &HashidsLinkEncoder{h}
}

func (e HashidsLinkEncoder) Encode(id int64) string {
	hash, err := e.h.EncodeInt64([]int64{id})
	if err != nil {
		panic(err)
	}

	return hash
}

func (e HashidsLinkEncoder) Decode(hash string) (int64, error) {
	ids, err := e.h.DecodeInt64WithError(hash)
	if err != nil {
		return 0, err
	}

	if len(ids) != 1 {
		return 0, errors.New("invalid hash")
	}

	return ids[0], nil
}
