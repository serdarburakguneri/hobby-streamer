package asset

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type AssetPage struct {
	Items            []Asset
	LastEvaluatedKey map[string]types.AttributeValue
}

func BuildPaginatedResponse(page *AssetPage) map[string]interface{} {
	return map[string]interface{}{
		"items":   page.Items,
		"nextKey": "", // Use shared.EncodeLastEvaluatedKey(page.LastEvaluatedKey) in handler
	}
}
