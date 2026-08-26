package entity

import "time"

type Repository struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Owner           Owner     `json:"owner"`
	HTMLURL         string    `json:"html_url"`
	Description     *string   `json:"description"` // Açıklama null gelebilir (*string)
	Fork            bool      `json:"fork"`
	Language        *string   `json:"language"` // JSON'da "language": null olabildiği için pointer
	Topics          []string  `json:"topics"`
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	WatchersCount   int       `json:"watchers_count"`
	OpenIssuesCount int       `json:"open_issues_count"`
	Size            int       `json:"size"`
	DefaultBranch   string    `json:"default_branch"`
	Visibility      string    `json:"visibility"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	PushedAt        time.Time `json:"pushed_at"`
}
