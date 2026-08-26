package repository

import (
	"context"
	"github-search-repo/internal/entity"
)

type ElasticRepository interface {
	BulkIndex(ctx context.Context, repos []entity.Repository) error
	Search(ctx context.Context, query string, language string, minStars int, page, size int) ([]entity.Repository, int64, error)
	GetExistingIDs(ctx context.Context, ids []int64) (map[int64]bool, error)
	Suggest(ctx context.Context, prefix string, size int) ([]string, error)
	GetMaxRepositoryID(ctx context.Context) (int64, error)
}
