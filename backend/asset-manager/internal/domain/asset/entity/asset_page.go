package entity

type AssetPage struct {
	Items   []*Asset
	HasMore bool
	LastKey map[string]interface{}
}
