package cdn

import "strings"

type Service interface {
	BuildPlayURL(key string) (cdnPrefix string, url string)
}

type service struct {
	prefix string
}

// Example prefix: "https://cdn.example.com"
func NewService(prefix string) Service {
	return &service{prefix: strings.TrimRight(prefix, "/")}
}

func (s *service) BuildPlayURL(key string) (string, string) {
	cleanKey := strings.TrimLeft(key, "/")
	if cleanKey == "" {
		return s.prefix, s.prefix
	}
	return s.prefix, s.prefix + "/" + cleanKey
}
