package main

import (
	"context"
	"github-search-repo/internal/controller"
	"github-search-repo/internal/repository"
	"github-search-repo/internal/service"
	"github-search-repo/internal/telemetry"
	"github.com/elastic/go-elasticsearch/v9"
	"log"
	"net/http"
	"os"
)

func main() {

	tp, err := telemetry.InitTracer("github-search-repo", context.Background())
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = tp.Shutdown(context.Background())
	}()
	es, err := elasticsearch.NewTyped(
		elasticsearch.WithAddresses("http://localhost:9200/"))

	if err != nil {
		panic(err)
	}
	indexName := "repositories"
	elasticRepo := repository.NewElasticRepository(es, indexName)

	githubToken := os.Getenv("GITHUB_TOKEN")
	syncService := service.NewSyncService(elasticRepo, githubToken)
	elasticService := service.NewElasticService(elasticRepo, syncService)

	repoController := controller.NewRepositoryController(elasticService)
	syncController := controller.NewSyncController(syncService)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/search", repoController.Search)
	mux.HandleFunc("/api/v1/sync", repoController.SyncHandler)
	mux.HandleFunc("/api/v1/suggest", repoController.SuggestHandler) // Yeni rota
	mux.HandleFunc("/api/v1/streamAllRepositories", syncController.StreamAllRepositories)
	port := ":8082"
	log.Printf("HTTP sunucusu %s portunda dinleniyor...", port)
	if err := http.ListenAndServe(port, enableCORS(mux)); err != nil {
		log.Fatalf("Sunucu hatası: %v", err)
	}
}
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
