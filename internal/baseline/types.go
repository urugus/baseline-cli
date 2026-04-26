package baseline

import "encoding/json"

type PageResponse[T any] struct {
	Data    []T `json:"data"`
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
	Total   int `json:"total"`
}

type EntityResponse[T any] struct {
	Data T `json:"data"`
}

type Vulnerability struct {
	ID        string          `json:"id"`
	NB        string          `json:"nb"`
	Title     string          `json:"title"`
	TitleJP   string          `json:"titleJp"`
	Severity  string          `json:"severity"`
	Status    string          `json:"status"`
	CustomID  string          `json:"customId"`
	Project   ProjectRef      `json:"project"`
	Asset     AssetRef        `json:"asset"`
	CVEs      []string        `json:"cves"`
	Refs      []string        `json:"refs"`
	Raw       json.RawMessage `json:"-"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updateAt"`
}

type ProjectRef struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Organization json.RawMessage `json:"organization"`
}

type AssetRef struct {
	ID          string `json:"id"`
	Project     string `json:"project"`
	Origin      string `json:"origin"`
	Type        string `json:"type"`
	ExternalID  string `json:"externalId"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}
