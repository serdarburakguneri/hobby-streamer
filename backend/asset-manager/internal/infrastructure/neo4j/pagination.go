package neo4j

func toPageParams(limit *int, offset *int) (int, map[string]interface{}) {
	l := 10
	if limit != nil {
		l = *limit
	}
	o := 0
	if offset != nil {
		o = *offset
	}
	return l, map[string]interface{}{"offset": o}
}
