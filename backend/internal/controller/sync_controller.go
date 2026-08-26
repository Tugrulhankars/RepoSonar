package controller

import (
	"github-search-repo/internal/service"
	"net/http"
)

type SyncController struct {
	syncService service.SyncService
}

func NewSyncController(syncService service.SyncService) *SyncController {
	return &SyncController{syncService: syncService}
}

func (s *SyncController) StreamAllRepositories(w http.ResponseWriter, r *http.Request) {

	err := s.syncService.StreamAllRepositories(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
