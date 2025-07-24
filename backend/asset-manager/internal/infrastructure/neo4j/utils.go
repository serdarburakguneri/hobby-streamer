package neo4j

import "time"

func ParseTimeToVO(s interface{}, voFunc func(time.Time) interface{}) interface{} {
	str, ok := s.(string)
	if !ok || str == "" {
		return voFunc(time.Now().UTC())
	}
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return voFunc(time.Now().UTC())
	}
	return voFunc(t)
}
