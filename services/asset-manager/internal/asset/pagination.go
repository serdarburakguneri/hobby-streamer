package asset

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type AssetPage struct {
	Items            []Asset
	LastEvaluatedKey map[string]types.AttributeValue
}

func DecodeLastEvaluatedKey(token string) (map[string]map[string]string, error) {
	raw, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	var decoded map[string]map[string]string
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}

func ToDynamoKey(encoded map[string]map[string]string) (map[string]types.AttributeValue, error) {
	if encoded == nil {
		return nil, nil
	}
	key := make(map[string]types.AttributeValue)
	for k, v := range encoded {
		switch v["type"] {
		case "S":
			key[k] = &types.AttributeValueMemberS{Value: v["value"]}
		case "N":
			key[k] = &types.AttributeValueMemberN{Value: v["value"]}
		default:
			return nil, errors.New("unsupported attribute type")
		}
	}
	return key, nil
}

func EncodeLastEvaluatedKey(raw map[string]types.AttributeValue) string {
	if raw == nil {
		return ""
	}
	enc := make(map[string]map[string]string)
	for k, v := range raw {
		switch t := v.(type) {
		case *types.AttributeValueMemberS:
			enc[k] = map[string]string{"type": "S", "value": t.Value}
		case *types.AttributeValueMemberN:
			enc[k] = map[string]string{"type": "N", "value": t.Value}
		}
	}
	b, _ := json.Marshal(enc)
	return base64.StdEncoding.EncodeToString(b)
}

func BuildPaginatedResponse(page *AssetPage) map[string]interface{} {
	return map[string]interface{}{
		"items":   page.Items,
		"nextKey": EncodeLastEvaluatedKey(page.LastEvaluatedKey),
	}
}