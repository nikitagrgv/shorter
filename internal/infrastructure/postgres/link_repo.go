package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nikitagrgv/shorter/internal/domain"
)

type LinkRepoPostgres struct {
	pool *pgxpool.Pool
}

func NewLinkRepoPostgres(pool *pgxpool.Pool) *LinkRepoPostgres {
	return &LinkRepoPostgres{pool: pool}
}

func (r LinkRepoPostgres) Create(ctx context.Context, link domain.Link, tokHash string) error {
	query := `
INSERT INTO links (id, short, long_url, created_at, access_token_hash)
VALUES ($1, $2, $3, $4, $5)
`
	_, err := r.pool.Exec(ctx, query,
		link.ID,
		link.Short,
		link.LongURL,
		link.CreatedAt,
		tokHash)

	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				switch pgErr.ConstraintName {
				case "links_pkey":
					return domain.ErrIdCollision
				case "links_short_unique":
					return domain.ErrShortCodeCollision
				case "links_access_token_hash_unique":
					return domain.ErrTokenCollision
				}
			}
		}
	}

	return err
}

func (r LinkRepoPostgres) FindByShortCode(ctx context.Context, code string) (domain.Link, error) {
	query := `
SELECT id, short, long_url, created_at
FROM links
WHERE short = $1
`

	var l domain.Link
	err := r.pool.QueryRow(ctx, query, code).Scan(&l.ID, &l.Short, &l.LongURL, &l.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Link{}, domain.ErrLinkNotFound
		}
		return domain.Link{}, err
	}

	return l, nil
}

func (r LinkRepoPostgres) FindByAccessTokenHash(ctx context.Context, tokHash string) (domain.Link, error) {
	query := `
SELECT id, short, long_url, created_at
FROM links
WHERE access_token_hash = $1
`

	var l domain.Link
	err := r.pool.QueryRow(ctx, query, tokHash).Scan(&l.ID, &l.Short, &l.LongURL, &l.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Link{}, domain.ErrLinkNotFound
		}
		return domain.Link{}, err
	}

	return l, nil
}
