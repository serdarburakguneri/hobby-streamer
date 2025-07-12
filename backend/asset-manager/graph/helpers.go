package graph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (r *mutationResolver) deleteAssetFiles(ctx context.Context, assetID string, files []struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}) error {
	type DeleteFilesRequest struct {
		AssetID string `json:"assetId"`
		Files   []struct {
			Bucket string `json:"bucket"`
			Key    string `json:"key"`
		} `json:"files"`
	}

	request := DeleteFilesRequest{
		AssetID: assetID,
		Files:   files,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal delete request: %w", err)
	}

	resp, err := http.Post(
		"http://localhost:4566/2015-03-31/functions/delete-asset-files/invocations",
		"application/json",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to call delete-asset-files function: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete-asset-files function returned status: %d", resp.StatusCode)
	}

	return nil
}
