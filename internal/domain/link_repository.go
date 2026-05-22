package domain

import "context"

type LinkRepo interface {
	Create(ctx context.Context, link Link, tokHash string) error
	FindByShortCode(ctx context.Context, code string) (Link, error)
	FindByAccessTokenHash(ctx context.Context, tokHash string) (Link, error)
}
