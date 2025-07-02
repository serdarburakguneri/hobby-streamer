package bucket

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type BucketPage struct {
	Items            []Bucket
	LastEvaluatedKey map[string]types.AttributeValue
}