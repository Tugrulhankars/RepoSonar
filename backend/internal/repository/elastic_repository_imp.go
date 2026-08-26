package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github-search-repo/internal/entity"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/sortorder"
)

type esRepo struct {
	client *elasticsearch.TypedClient
	index  string
}

func (e *esRepo) Suggest(ctx context.Context, prefix string, size int) ([]string, error) {
	if size <= 0 || size > 10 {
		size = 5
	}

	boolQuery := types.NewBoolQuery()
	boolQuery.Should = append(boolQuery.Should,
		types.Query{
			Prefix: map[string]types.PrefixQuery{
				"name.keyword": {Value: prefix},
			},
		},
		types.Query{
			MatchPhrasePrefix: map[string]types.MatchPhrasePrefixQuery{
				"name": {Query: prefix},
			},
		},
	)

	// Sıralama ve arama isteği
	res, err := e.client.Search().
		Index(e.index).
		Header("Content-Type", "application/json").
		Header("Accept", "application/json").
		Size(size).
		Query(&types.Query{Bool: boolQuery}).
		Sort(&types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				"stargazers_count": {
					Order: &sortorder.Desc,
				},
			},
		}).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("suggest sorgu hatası: %w", err)
	}

	var suggestions []string
	for _, hit := range res.Hits.Hits {
		var repo entity.Repository
		if err := json.Unmarshal(hit.Source_, &repo); err == nil {
			suggestions = append(suggestions, repo.FullName)
		}
	}

	return suggestions, nil
}

func NewElasticRepository(client *elasticsearch.TypedClient, index string) ElasticRepository {
	return &esRepo{
		client: client,
		index:  index,
	}
}

// BulkIndex: GitHub'dan gelen verileri Elasticsearch'e toplu olarak kaydeder / günceller
func (e *esRepo) BulkIndex(ctx context.Context, repos []entity.Repository) error {
	var buf bytes.Buffer

	for _, repo := range repos {
		meta := fmt.Sprintf(`{ "index" : { "_index" : "%s", "_id" : "%d" } }%s`, e.index, repo.ID, "\n")
		data, err := json.Marshal(repo)
		if err != nil {
			return err
		}

		buf.WriteString(meta)
		buf.Write(data)
		buf.WriteString("\n")
	}

	// TypedClient ile Raw NDJSON Bulk işlemi
	res, err := e.client.Bulk().
		Header("Content-Type", "application/x-ndjson").
		Header("Accept", "application/json").
		Raw(&buf).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("bulk indexing isteği başarısız: %w", err)
	}

	if res.Errors {
		return fmt.Errorf("bulk indexing sırasında bazı kayıtlar eklenemedi")
	}

	return nil
}

// Search: Full-text ve filtreli arama sorgusunu çalıştırır
func (e *esRepo) Search(ctx context.Context, query string, language string, minStars int, page, size int) ([]entity.Repository, int64, error) {
	from := (page - 1) * size
	boolQuery := types.NewBoolQuery()

	// 1. Full-text Arama
	if query != "" {
		boolQuery.Must = append(boolQuery.Must, types.Query{
			MultiMatch: &types.MultiMatchQuery{
				Query:     query,
				Fields:    []string{"name^3", "description^1", "topics^2"},
				Fuzziness: "AUTO",
			},
		})
	} else {
		boolQuery.Must = append(boolQuery.Must, types.Query{
			MatchAll: &types.MatchAllQuery{},
		})
	}

	// 2. Dil Filtresi (TermQuery)
	if language != "" {
		caseInsensitive := true
		boolQuery.Filter = append(boolQuery.Filter, types.Query{
			Term: map[string]types.TermQuery{
				"language": {
					Value:           language,
					CaseInsensitive: &caseInsensitive,
				},
			},
		})
	}

	// 3. Yıldız Sayısı Filtresi
	if minStars > 0 {
		minStarsFloat := types.Float64(minStars)
		boolQuery.Filter = append(boolQuery.Filter, types.Query{
			Range: map[string]types.RangeQuery{
				"stargazers_count": types.NumberRangeQuery{
					Gte: &minStarsFloat,
				},
			},
		})
	}
	// 4. Elasticsearch Search Request
	res, err := e.client.Search().
		Index(e.index).
		Header("Content-Type", "application/json").
		Header("Accept", "application/json").
		Request(&search.Request{
			From:  &from,
			Size:  &size,
			Query: &types.Query{Bool: boolQuery},
			Sort: []types.SortCombinations{
				&types.SortOptions{
					SortOptions: map[string]types.FieldSort{
						"stargazers_count": {
							Order: &sortorder.Desc,
						},
					},
				},
			},
		}).
		Do(ctx)

	if err != nil {
		// Hatanın detayını terminale basarak inceleyelim
		log.Printf("❌ Elasticsearch Search Hatası: %v", err)
		return nil, 0, fmt.Errorf("elastic arama hatası: %w", err)
	}

	// 5. Sonuçları Diziye Aktar
	repos := make([]entity.Repository, 0, len(res.Hits.Hits))
	for _, hit := range res.Hits.Hits {
		var repo entity.Repository
		if err := json.Unmarshal(hit.Source_, &repo); err != nil {
			continue
		}
		repos = append(repos, repo)
	}

	var totalCount int64
	if res.Hits.Total != nil {
		totalCount = res.Hits.Total.Value
	}

	return repos, totalCount, nil
}

// GetExistingIDs: Elasticsearch'te belirtilen ID'lerden var olanları sorgular
func (e *esRepo) GetExistingIDs(ctx context.Context, ids []int64) (map[int64]bool, error) {
	if len(ids) == 0 {
		return make(map[int64]bool), nil
	}

	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = strconv.FormatInt(id, 10)
	}

	size := len(idStrings)
	res, err := e.client.Search().
		Index(e.index).
		Header("Content-Type", "application/json").
		Header("Accept", "application/json").
		Request(&search.Request{
			Size: &size,
			Query: &types.Query{
				Ids: &types.IdsQuery{
					Values: idStrings,
				},
			},
			Source_: &types.SourceFilter{
				Includes: []string{"id"},
			},
		}).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("mevcut ID'leri sorgulama hatası: %w", err)
	}

	existingMap := make(map[int64]bool)
	for _, hit := range res.Hits.Hits {
		if hit.Id_ != nil {
			docID, err := strconv.ParseInt(*hit.Id_, 10, 64)
			if err == nil {
				existingMap[docID] = true
			}
		}
	}

	return existingMap, nil
}

func (e *esRepo) GetMaxRepositoryID(ctx context.Context) (int64, error) {
	// ID alanına göre büyükten küçüğe 1 tane kayıt çek
	res, err := e.client.Search().
		Index(e.index).
		Size(1).
		Sort(&types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				"id": {Order: &sortorder.Desc},
			},
		}).
		Do(ctx)

	if err != nil || len(res.Hits.Hits) == 0 {
		return 0, err
	}

	var repo entity.Repository
	if err := json.Unmarshal(res.Hits.Hits[0].Source_, &repo); err != nil {
		return 0, err
	}

	return repo.ID, nil
}
