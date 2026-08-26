package service

import (
	"context"
	"fmt"
	"github-search-repo/internal/entity"
	"github-search-repo/internal/repository"
	"log"
)

type elasticService struct {
	repo        repository.ElasticRepository
	syncService SyncService
}

func NewElasticService(repo repository.ElasticRepository, syncService SyncService) ElasticService {
	return &elasticService{
		repo:        repo,
		syncService: syncService,
	}
}

func (e *elasticService) BulkIndex(ctx context.Context) error {
	languages := []string{"go", "rust", "java", "c#"}

	// 1. GitHub'dan verileri çek
	resp, err := e.syncService.SyncPopularRepositories(ctx, languages)
	if err != nil {
		return fmt.Errorf("github sync hatası: %w", err)
	}

	if len(resp) == 0 {
		return nil
	}

	// 2. Elasticsearch'te zaten var olan repo ID'lerini sorgula
	ids := make([]int64, len(resp))
	for i, repo := range resp {
		ids[i] = repo.ID
	}

	existingIDs, err := e.repo.GetExistingIDs(ctx, ids)
	if err != nil {
		log.Printf("Mevcut ID'ler kontrol edilirken uyarı (devam ediliyor): %v", err)
		existingIDs = make(map[int64]bool)
	}

	// 3. Sadece Elasticsearch'te var OLMAYAN yeni repoları filtrele
	var newRepos []entity.Repository
	for _, repo := range resp {
		if !existingIDs[repo.ID] {
			newRepos = append(newRepos, repo)
		}
	}

	if len(newRepos) == 0 {
		log.Println("ℹ️ Çekilen tüm repolar zaten Elasticsearch'te mevcut, yeni veri eklenmedi.")
		return nil
	}

	log.Printf("ℹ️ Çekilen %d repodan %d tanesi yeni. Elasticsearch'e ekleniyor...", len(resp), len(newRepos))

	// 4. Çekilen yeni verileri Elasticsearch'e yaz
	if err := e.repo.BulkIndex(ctx, newRepos); err != nil {
		return fmt.Errorf("elasticsearch bulk index hatası: %w", err)
	}

	return nil
}

func (e *elasticService) Search(ctx context.Context, query string, language string, minStars int, page, size int) ([]entity.Repository, int64, error) {
	// 1. Repo çağrısı 3 değer döner (repos, totalCount, err)
	repos, totalCount, err := e.repo.Search(ctx, query, language, minStars, page, size)
	if err != nil {
		return nil, 0, err
	}

	return repos, totalCount, nil
}
func (e *elasticService) Suggest(ctx context.Context, prefix string, size int) ([]string, error) {
	if prefix == "" {
		return []string{}, nil
	}
	return e.repo.Suggest(ctx, prefix, size)
}
func (e *elasticService) GetMaxRepositoryID(ctx context.Context) (int64, error) {

	id, err := e.repo.GetMaxRepositoryID(ctx)
	if err != nil {
		return -1, err
	}

	return id, err
}
