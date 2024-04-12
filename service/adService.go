package service

import (
	"context"

	"github.com/hunick1234/DcardBackend/model"
	"github.com/hunick1234/DcardBackend/model/ad"
)

type AdService interface {
	Aggregate(context.Context, model.Filter, any) error
	FindByFilter(context.Context, *ad.AdQuery) (*[]ad.ResponseAd, error)
	Store(context.Context, *ad.AD) error
}

type adService struct {
	repo model.Repository[ad.AD, ad.AdQuery, ad.ResponseAd]
}

func NewAdService(repo model.Repository[ad.AD, ad.AdQuery, ad.ResponseAd]) AdService {
	return &adService{repo: repo}
}

func (s *adService) FindByFilter(ctx context.Context, q *ad.AdQuery) (*[]ad.ResponseAd, error) {
	return s.repo.FindByFilter(ctx, q)
}

func (s *adService) Store(ctx context.Context, ad *ad.AD) error {
	return s.repo.Store(ctx, ad)
}

func (s *adService) Aggregate(ctx context.Context, filter model.Filter, results any) error {
	return s.repo.Aggregate(ctx, filter, results)
}
