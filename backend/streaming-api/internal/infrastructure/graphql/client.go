package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type Client struct {
	httpClient      *http.Client
	serviceClient   auth.ServiceClientInterface
	assetManagerURL string
	logger          *logger.Logger
}

func NewClient(serviceClient auth.ServiceClientInterface, assetManagerURL string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		serviceClient:   serviceClient,
		assetManagerURL: assetManagerURL,
		logger:          logger.Get().WithService("graphql-client"),
	}
}

func (c *Client) Query(ctx context.Context, query string, variables map[string]interface{}, result interface{}) error {
	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return errors.NewTransientError("failed to marshal request", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.assetManagerURL+"/graphql", bytes.NewBuffer(jsonBody))
	if err != nil {
		return errors.NewTransientError("failed to create request", err)
	}

	req.Header.Set("Content-Type", "application/json")

	token, err := c.serviceClient.GetServiceToken(ctx)
	if err != nil {
		return errors.NewTransientError("failed to get service token", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	c.logger.Debug("Sending GraphQL request", "url", req.URL.String(), "query", query)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.NewTransientError("failed to send request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.NewTransientError(fmt.Sprintf("GraphQL request failed with status: %d", resp.StatusCode), nil)
	}

	var graphQLResponse struct {
		Data   interface{} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&graphQLResponse); err != nil {
		return errors.NewTransientError("failed to decode response", err)
	}

	if len(graphQLResponse.Errors) > 0 {
		errorMessages := make([]string, len(graphQLResponse.Errors))
		for i, err := range graphQLResponse.Errors {
			errorMessages[i] = err.Message
		}
		return errors.NewTransientError(fmt.Sprintf("GraphQL errors: %v", errorMessages), nil)
	}

	if graphQLResponse.Data == nil {
		return errors.NewNotFoundError("no data returned from GraphQL", nil)
	}

	dataJSON, err := json.Marshal(graphQLResponse.Data)
	if err != nil {
		return errors.NewTransientError("failed to marshal response data", err)
	}

	if err := json.Unmarshal(dataJSON, result); err != nil {
		return errors.NewTransientError("failed to unmarshal response data", err)
	}

	return nil
}
