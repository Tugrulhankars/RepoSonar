package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github-search-repo/internal/dto/response"
	"github-search-repo/internal/entity"
	"github-search-repo/internal/repository"
	"go.opentelemetry.io/otel"
	"log"
	"net/http"
	"net/url"
	"time"
)

type SyncService interface {
	SyncPopularRepositories(ctx context.Context, languages []string) ([]entity.Repository, error)
	StreamAllRepositories(ctx context.Context) error
}

var tracer = otel.Tracer("Sync-Service")

type syncService struct {
	elasticRepo repository.ElasticRepository
	httpClient  *http.Client
	githubToken string
}

func NewSyncService(elasticRepo repository.ElasticRepository, githubToken string) SyncService {
	return &syncService{
		elasticRepo: elasticRepo,
		githubToken: githubToken,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (s *syncService) SyncPopularRepositories(ctx context.Context, languages []string) ([]entity.Repository, error) {
	var allRepositories []entity.Repository

	for _, lang := range languages {
		// Her bir istek için bağımsız 45 saniyelik timeout context'i oluşturulur
		reqCtx, cancel := context.WithTimeout(context.Background(), 45*time.Second)

		// 1. URL ve Query Hazırlığı
		rawQuery := fmt.Sprintf("language:%s stars:>500", lang)
		encodedQuery := url.QueryEscape(rawQuery)
		apiURL := fmt.Sprintf("https://api.github.com/search/repositories?q=%s&sort=stars&order=desc&per_page=100&page=1", encodedQuery)

		// 2. HTTP GET İsteği
		req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, apiURL, nil)
		if err != nil {
			cancel()
			log.Printf("[%s] İstek hatası: %v", lang, err)
			continue
		}

		req.Header.Set("Accept", "application/vnd.github.v3+json")
		req.Header.Set("User-Agent", "GitHub-Search-App/1.0")
		if s.githubToken != "" {
			req.Header.Set("Authorization", "Bearer "+s.githubToken)
		}

		resp, err := s.httpClient.Do(req)
		if err != nil {
			cancel()
			log.Printf("[%s] İstek başarısız: %v", lang, err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("[%s] GitHub Hatası (Status %d)", lang, resp.StatusCode)
			resp.Body.Close()
			cancel()
			continue
		}

		// 3. Response Parse
		var searchResult response.GitHubSearchResponse
		err = json.NewDecoder(resp.Body).Decode(&searchResult)
		resp.Body.Close()
		cancel()

		if err != nil {
			log.Printf("[%s] JSON Decode Hatası: %v", lang, err)
			continue
		}

		// 4. Çekilen repoları ana listeye ekle
		allRepositories = append(allRepositories, searchResult.Items...)

		log.Printf("✅ [%s] %d repo başarıyla listeye eklendi", lang, len(searchResult.Items))

		time.Sleep(2 * time.Second)
	}

	return allRepositories, nil
}

func (s *syncService) StreamAllRepositories(ctx context.Context) error {
	// 1. Elasticsearch'teki en son indekslenmiş en büyük ID'yi getir
	/*lastRepoID, err := s.elasticRepo.GetMaxRepositoryID(ctx)
	if err != nil {
		log.Printf("Mevcut ID okunamadı, 0'dan başlanıyor: %v", err)
		lastRepoID = 0
	} else {
		log.Printf("🔄 Senkronizasyon %d nolu ID'den devam ediyor...", lastRepoID)
	}*/
	lastRepoID := int64(0)
	for {
		ctx, span := tracer.Start(ctx, "SyncService.StreamAllRepositories")
		defer span.End()

		apiURL := fmt.Sprintf("https://api.github.com/repositories?since=%d", lastRepoID)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
		if err != nil {
			span.AddEvent(err.Error())
			return err
		}
		req.Header.Set("Authorization", "Bearer "+s.githubToken)
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err := s.httpClient.Do(req)
		if err != nil {
			return err
		}

		var repos []entity.Repository
		err = json.NewDecoder(resp.Body).Decode(&repos)
		resp.Body.Close()

		if err != nil || len(repos) == 0 {
			break
		}

		// 2. Elasticsearch'e kaydet
		if err := s.elasticRepo.BulkIndex(ctx, repos); err != nil {
			return fmt.Errorf("bulk index hatası: %w", err)
		}

		// 3. Son reponun ID'sini al
		lastRepoID = repos[len(repos)-1].ID

		time.Sleep(1 * time.Second)
	}
	return nil
}
