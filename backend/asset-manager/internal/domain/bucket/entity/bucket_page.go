package entity

type BucketPage struct {
	Items   []*Bucket
	HasMore bool
	LastKey map[string]interface{}
}
