package graphql

import (
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
)

func TestGraphQLSchema(t *testing.T) {
	resolver := &Resolver{}
	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))

	c := client.New(srv)

	t.Run("introspection query", func(t *testing.T) {
		var resp struct {
			Schema struct {
				QueryType struct {
					Name string `json:"name"`
				} `json:"queryType"`
				MutationType struct {
					Name string `json:"name"`
				} `json:"mutationType"`
			} `json:"__schema"`
		}

		err := c.Post(`
			query IntrospectionQuery {
				__schema {
					queryType {
						name
					}
					mutationType {
						name
					}
				}
			}
		`, &resp)

		if err != nil {
			t.Errorf("Introspection query failed: %v", err)
		}

		if resp.Schema.QueryType.Name != "Query" {
			t.Errorf("Expected Query type, got %s", resp.Schema.QueryType.Name)
		}

		if resp.Schema.MutationType.Name != "Mutation" {
			t.Errorf("Expected Mutation type, got %s", resp.Schema.MutationType.Name)
		}
	})

	t.Run("bucket type exists", func(t *testing.T) {
		var resp struct {
			Schema struct {
				Types []struct {
					Name string `json:"name"`
				} `json:"types"`
			} `json:"__schema"`
		}

		err := c.Post(`
			query IntrospectionQuery {
				__schema {
					types {
						name
					}
				}
			}
		`, &resp)

		if err != nil {
			t.Errorf("Introspection query failed: %v", err)
		}

		bucketTypeExists := false
		for _, t := range resp.Schema.Types {
			if t.Name == "Bucket" {
				bucketTypeExists = true
				break
			}
		}

		if !bucketTypeExists {
			t.Error("Bucket type not found in schema")
		}
	})
}
