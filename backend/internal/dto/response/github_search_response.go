package response

import "github-search-repo/internal/entity"

type GitHubSearchResponse struct {
	TotalCount        int                 `json:"total_count"`
	IncompleteResults bool                `json:"incomplete_results"`
	Items             []entity.Repository `json:"items"`
}
