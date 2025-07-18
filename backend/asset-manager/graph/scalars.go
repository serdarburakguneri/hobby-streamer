package graph

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// MarshalJSONObject marshals JSONObject to GraphQL
func MarshalJSONObject(v map[string]interface{}) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		data, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		if _, err := w.Write(data); err != nil {
			panic(err)
		}
	})
}

// UnmarshalJSONObject unmarshals JSONObject from GraphQL
func UnmarshalJSONObject(v interface{}) (map[string]interface{}, error) {
	switch v := v.(type) {
	case map[string]interface{}:
		return v, nil
	case string:
		var result map[string]interface{}
		err := json.Unmarshal([]byte(v), &result)
		return result, err
	case []byte:
		var result map[string]interface{}
		err := json.Unmarshal(v, &result)
		return result, err
	default:
		return nil, fmt.Errorf("JSONObject must be a map[string]interface{}, string, or []byte")
	}
}

// MarshalTime marshals Time to GraphQL
func MarshalTime(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if _, err := io.WriteString(w, fmt.Sprintf("%q", t.Format(time.RFC3339))); err != nil {
			panic(err)
		}
	})
}

// UnmarshalTime unmarshals Time from GraphQL
func UnmarshalTime(v interface{}) (time.Time, error) {
	switch v := v.(type) {
	case string:
		return time.Parse(time.RFC3339, v)
	case time.Time:
		return v, nil
	default:
		return time.Time{}, fmt.Errorf("Time must be a string or time.Time")
	}
}
