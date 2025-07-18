package model

type BucketsResponse struct {
	Buckets []Bucket `json:"buckets"`
	Count   int      `json:"count"`
}

type AssetsResponse struct {
	Assets []Asset `json:"assets"`
	Count  int     `json:"count"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
