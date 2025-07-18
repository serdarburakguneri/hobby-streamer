package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type GraphQLClient struct {
	logger        *logger.Logger
	serviceClient auth.ServiceClientInterface
}

func NewGraphQLClient(serviceClient auth.ServiceClientInterface) *GraphQLClient {
	return &GraphQLClient{
		logger:        logger.WithService("graphql-client"),
		serviceClient: serviceClient,
	}
}

func (g *GraphQLClient) ExecuteQuery(ctx context.Context, url, query string, response interface{}) error {
	requestBody := map[string]string{
		"query": query,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return errors.NewInternalError("failed to marshal request body", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return errors.NewInternalError("failed to create request", err)
	}

	req.Header.Set("Content-Type", "application/json")

	g.logger.Info("Getting service token for request", "url", url)
	authHeader, err := g.serviceClient.GetAuthorizationHeader(ctx)
	if err != nil {
		g.logger.WithError(err).Error("Failed to get service token")
		return errors.NewExternalError("failed to get service token", err)
	}

	g.logger.Info("Service token obtained successfully", "auth_header_length", len(authHeader))
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return errors.NewTransientError("failed to make request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		g.logger.Error("GraphQL request failed", "status_code", resp.StatusCode, "url", url)
		return errors.NewExternalError(fmt.Sprintf("unexpected status code: %d", resp.StatusCode), nil)
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return errors.NewInternalError("failed to decode response", err)
	}

	return nil
}

func (g *GraphQLClient) ExecuteQueryWithCircuitBreaker(ctx context.Context, circuitBreaker *errors.CircuitBreaker, url, query string, response interface{}) error {
	return circuitBreaker.Execute(ctx, func() error {
		return g.ExecuteQuery(ctx, url, query, response)
	})
}

func (g *GraphQLClient) HandleGraphQLErrors(response interface{}) error {
	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return nil
	}

	if graphQLErrors, exists := responseMap["errors"]; exists {
		if errorList, ok := graphQLErrors.([]interface{}); ok && len(errorList) > 0 {
			return errors.NewExternalError(fmt.Sprintf("GraphQL errors: %v", graphQLErrors), nil)
		}
	}

	return nil
}
