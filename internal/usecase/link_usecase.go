package usecase

import (
	"context"

	"github.com/nikitagrgv/shorter/internal/domain"
)

type LinkUsecase struct {
	tp          domain.TimeProvider
	repo        domain.LinkRepo
	idGen       domain.IdGenerator
	shortGen    domain.ShortLinkEncoder
	tokenGen    domain.TokenGenerator
	tokenHasher domain.TokenHasher
}

func NewLinkUsecase(
	tp domain.TimeProvider,
	repo domain.LinkRepo,
	idGen domain.IdGenerator,
	shortLinkEncoder domain.ShortLinkEncoder,
	tokenGen domain.TokenGenerator,
	tokenHasher domain.TokenHasher) *LinkUsecase {
	return &LinkUsecase{tp, repo, idGen, shortLinkEncoder, tokenGen, tokenHasher}
}

type CreateLinkCommand struct {
	LongURL string
}

type CreateLinkResult struct {
	Link        domain.Link
	AccessToken string // secret for admin actions: see clicks count, delete, etc
}

func (u *LinkUsecase) CreateLink(ctx context.Context, command CreateLinkCommand) (CreateLinkResult, error) {
	now := u.tp.NowUTC()
	id := u.idGen.Generate()
	link := domain.Link{
		ID:        id,
		Short:     u.shortGen.Encode(id),
		LongURL:   command.LongURL,
		CreatedAt: now,
	}

	tok := u.tokenGen.Generate()
	tokHash := u.tokenHasher.Hash(tok)

	err := u.repo.Create(ctx, link, tokHash)
	if err != nil {
		return CreateLinkResult{}, err
	}

	return CreateLinkResult{
		Link:        link,
		AccessToken: tok,
	}, nil
}

func (u *LinkUsecase) GetLinkByToken(ctx context.Context, token string) (domain.Link, error) {
	tokHash := u.tokenHasher.Hash(token)
	link, err := u.repo.FindByAccessTokenHash(ctx, tokHash)
	if err != nil {
		return domain.Link{}, err
	}
	return link, nil
}

func (u *LinkUsecase) GetLinkByCode(ctx context.Context, code string) (domain.Link, error) {
	link, err := u.repo.FindByShortCode(ctx, code)
	if err != nil {
		return domain.Link{}, err
	}
	return link, nil
}
