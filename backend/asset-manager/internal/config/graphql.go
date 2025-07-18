package config

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/graph"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/bucket"
)

type GraphQLConfig struct {
	Handler *handler.Server
	Router  *mux.Router
}

func NewGraphQLConfig(assetService *asset.Service, bucketService *bucket.Service) *GraphQLConfig {
	resolver := graph.NewResolver(assetService, bucketService)
	schema := graph.NewExecutableSchema(graph.Config{Resolvers: resolver})
	gqlHandler := handler.New(schema)

	gqlHandler.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				allowedOrigins := []string{"http://localhost:8081", "http://localhost:3000", "http://localhost:8080"}
				for _, allowed := range allowedOrigins {
					if origin == allowed {
						return true
					}
				}
				return false
			},
		},
		KeepAlivePingInterval: 10 * time.Second,
	})

	gqlHandler.AddTransport(&transport.Options{})
	gqlHandler.AddTransport(&transport.GET{})
	gqlHandler.AddTransport(&transport.POST{})
	gqlHandler.AddTransport(&transport.MultipartForm{})
	gqlHandler.Use(extension.Introspection{})
	gqlHandler.Use(extension.FixedComplexityLimit(1000))

	router := mux.NewRouter()
	router.Handle("/graphql", gqlHandler).Methods("GET", "POST", "OPTIONS")
	router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))

	return &GraphQLConfig{
		Handler: gqlHandler,
		Router:  router,
	}
}
