package queries

type GetAssetQuery struct {
	ID   string `json:"id,omitempty"`
	Slug string `json:"slug,omitempty"`
}

type ListAssetsQuery struct {
	Limit  *int `json:"limit"`
	Offset *int `json:"offset"`
}

type SearchAssetsQuery struct {
	Query  string `json:"query"`
	Limit  *int   `json:"limit"`
	Offset *int   `json:"offset"`
}
