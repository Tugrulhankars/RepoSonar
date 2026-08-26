package controller

import (
	"encoding/json"
	"github-search-repo/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type RepositoryController struct {
	elasticService service.ElasticService
}

func NewRepositoryController(elasticService service.ElasticService) *RepositoryController {
	return &RepositoryController{
		elasticService: elasticService,
	}
}

func (s *RepositoryController) Search(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Yalnızca GET metodu desteklenir"})
		return
	}

	queryParams := r.URL.Query()

	// query parametrelerini esnek parse etme
	query := queryParams.Get("q")

	language := queryParams.Get("lang")
	if language == "" {
		language = queryParams.Get("language")
	}

	minStarsStr := queryParams.Get("min_stars")
	if minStarsStr == "" {
		minStarsStr = queryParams.Get("minStars")
	}
	if minStarsStr == "" {
		minStarsStr = queryParams.Get("stars")
	}
	if minStarsStr == "" {
		minStarsStr = queryParams.Get("min_star")
	}
	if minStarsStr == "" {
		minStarsStr = queryParams.Get("minstars")
	}

	minStarsStr = strings.ReplaceAll(minStarsStr, ",", "")
	minStarsStr = strings.TrimSpace(minStarsStr)
	minStars, _ := strconv.Atoi(minStarsStr)
	page, _ := strconv.Atoi(queryParams.Get("page"))
	if page < 1 {
		page = 1
	}

	size, _ := strconv.Atoi(queryParams.Get("size"))
	if size < 1 || size > 100 {
		size = 10
	}

	// Service çağrısı
	repos, totalCount, err := s.elasticService.Search(
		r.Context(),
		query,
		language,
		minStars,
		page,
		size,
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Arama sırasında hata oluştu: " + err.Error(),
		})
		return
	}

	// JSON Yanıtı
	response := map[string]interface{}{
		"total_hits": totalCount,
		"page":       page,
		"size":       size,
		"data":       repos,
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *RepositoryController) SyncHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Yalnızca POST metodu desteklenir"})
		return
	}

	err := s.elasticService.BulkIndex(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Senkronizasyon hatası: " + err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "GitHub repository senkronizasyonu başarıyla tamamlandı",
	})
}
func (ctrl *RepositoryController) SuggestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Yalnızca GET metodu desteklenir"})
		return
	}

	prefix := r.URL.Query().Get("prefix")
	if prefix == "" {
		writeJSON(w, http.StatusOK, map[string]interface{}{"suggestions": []string{}})
		return
	}

	suggestions, err := ctrl.elasticService.Suggest(r.Context(), prefix, 5)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"prefix":      prefix,
		"suggestions": suggestions,
	})
}

// JSON formatında response basmak için yardımcı fonksiyon
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
