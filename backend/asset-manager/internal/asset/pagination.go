package asset

type AssetPage struct {
	Items []Asset
}

func BuildPaginatedResponse(page *AssetPage) map[string]interface{} {
	return map[string]interface{}{
		"items": page.Items,
	}
}
