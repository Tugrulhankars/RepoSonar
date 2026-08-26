package service

import (
	"context"
	"github-search-repo/internal/entity"
)

type ElasticService interface {
	BulkIndex(ctx context.Context) error
	Search(ctx context.Context, query string, language string, minStars int, page, size int) ([]entity.Repository, int64, error)
	Suggest(ctx context.Context, prefix string, size int) ([]string, error)
	GetMaxRepositoryID(ctx context.Context) (int64, error)
}
