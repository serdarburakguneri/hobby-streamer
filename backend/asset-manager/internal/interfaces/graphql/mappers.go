package graphql

import (
	assetCommands "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/commands"
	bucketCommands "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/bucket/commands"

	assetvo "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
	bucketvo "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket/valueobjects"
)

func MapCreateAssetInput(input CreateAssetInput) (assetCommands.CreateAssetCommand, error) {
	slug, err := assetvo.NewSlug(input.Slug)
	if err != nil {
		return assetCommands.CreateAssetCommand{}, err
	}
	var title *assetvo.Title
	if input.Title != nil {
		t, err := assetvo.NewTitle(*input.Title)
		if err != nil {
			return assetCommands.CreateAssetCommand{}, err
		}
		title = t
	}
	var at *assetvo.AssetType
	if input.Type != nil {
		t, err := assetvo.NewAssetType(*input.Type)
		if err != nil {
			return assetCommands.CreateAssetCommand{}, err
		}
		at = t
	}
	var owner *assetvo.OwnerID
	if input.OwnerID != nil {
		o, err := assetvo.NewOwnerID(*input.OwnerID)
		if err != nil {
			return assetCommands.CreateAssetCommand{}, err
		}
		owner = o
	}
	var parent *assetvo.AssetID
	if input.ParentID != nil {
		p, err := assetvo.NewAssetID(*input.ParentID)
		if err != nil {
			return assetCommands.CreateAssetCommand{}, err
		}
		parent = p
	}
	var genre *assetvo.Genre
	if input.Genre != nil {
		g, err := assetvo.NewGenre(*input.Genre)
		if err != nil {
			return assetCommands.CreateAssetCommand{}, err
		}
		genre = g
	}
	var genres *assetvo.Genres
	if input.Genres != nil && len(input.Genres) > 0 {
		gs, err := assetvo.NewGenres(input.Genres)
		if err != nil {
			return assetCommands.CreateAssetCommand{}, err
		}
		genres = gs
	}
	var tags *assetvo.Tags
	if input.Tags != nil && len(input.Tags) > 0 {
		ts, err := assetvo.NewTags(input.Tags)
		if err != nil {
			return assetCommands.CreateAssetCommand{}, err
		}
		tags = ts
	}
	return assetCommands.CreateAssetCommand{
		Slug:      *slug,
		Title:     title,
		AssetType: at,
		Genre:     genre,
		Genres:    genres,
		Tags:      tags,
		OwnerID:   owner,
		ParentID:  parent,
	}, nil
}

func MapDeleteAssetInput(id string) (assetCommands.DeleteAssetCommand, error) {
	idVO, err := assetvo.NewAssetID(id)
	if err != nil {
		return assetCommands.DeleteAssetCommand{}, err
	}
	return assetCommands.DeleteAssetCommand{ID: *idVO}, nil
}

func MapCreateBucketInput(input BucketInput) (bucketCommands.CreateBucketCommand, error) {
	var owner *bucketvo.OwnerID
	if input.OwnerID != nil {
		o, err := bucketvo.NewOwnerID(*input.OwnerID)
		if err != nil {
			return bucketCommands.CreateBucketCommand{}, err
		}
		owner = o
	}
	var desc *bucketvo.BucketDescription
	if input.Description != nil {
		d, err := bucketvo.NewBucketDescription(*input.Description)
		if err != nil {
			return bucketCommands.CreateBucketCommand{}, err
		}
		desc = d
	}
	var typ *bucketvo.BucketType
	if input.Type != nil {
		t, err := bucketvo.NewBucketType(*input.Type)
		if err != nil {
			return bucketCommands.CreateBucketCommand{}, err
		}
		typ = t
	}
	var stat *bucketvo.BucketStatus
	if input.Status != nil {
		s, err := bucketvo.NewBucketStatus(*input.Status)
		if err != nil {
			return bucketCommands.CreateBucketCommand{}, err
		}
		stat = s
	}
	return bucketCommands.CreateBucketCommand{Name: *input.Name, Key: *input.Key, OwnerID: owner, Description: desc, Type: typ, Status: stat, Metadata: map[string]interface{}{}}, nil
}

func MapUpdateBucketInput(id string, input BucketInput) (bucketCommands.UpdateBucketCommand, error) {
	idVO, err := bucketvo.NewBucketID(id)
	if err != nil {
		return bucketCommands.UpdateBucketCommand{}, err
	}
	cmd := bucketCommands.UpdateBucketCommand{ID: *idVO}
	if input.Name != nil {
		n, err := bucketvo.NewBucketName(*input.Name)
		if err != nil {
			return cmd, err
		}
		cmd.Name = n
	}
	if input.Description != nil {
		d, err := bucketvo.NewBucketDescription(*input.Description)
		if err != nil {
			return cmd, err
		}
		cmd.Description = d
	}
	if input.Status != nil {
		s, err := bucketvo.NewBucketStatus(*input.Status)
		if err != nil {
			return cmd, err
		}
		cmd.Status = s
	}
	if input.Type != nil {
		t, err := bucketvo.NewBucketType(*input.Type)
		if err != nil {
			return cmd, err
		}
		cmd.Type = t
	}
	if input.OwnerID != nil {
		o, err := bucketvo.NewOwnerID(*input.OwnerID)
		if err != nil {
			return cmd, err
		}
		cmd.OwnerID = o
	}
	return cmd, nil
}

func MapDeleteBucketInput(id string) (bucketCommands.DeleteBucketCommand, error) {
	idVO, err := bucketvo.NewBucketID(id)
	if err != nil {
		return bucketCommands.DeleteBucketCommand{}, err
	}
	return bucketCommands.DeleteBucketCommand{ID: *idVO}, nil
}

func MapAddAssetToBucketInput(input AddAssetToBucketInput) (bucketCommands.AddAssetToBucketCommand, error) {
	idVO, err := bucketvo.NewBucketID(input.BucketID)
	if err != nil {
		return bucketCommands.AddAssetToBucketCommand{}, err
	}
	return bucketCommands.AddAssetToBucketCommand{BucketID: *idVO, AssetID: input.AssetID}, nil
}

func MapRemoveAssetFromBucketInput(input RemoveAssetFromBucketInput) (bucketCommands.RemoveAssetFromBucketCommand, error) {
	idVO, err := bucketvo.NewBucketID(input.BucketID)
	if err != nil {
		return bucketCommands.RemoveAssetFromBucketCommand{}, err
	}
	return bucketCommands.RemoveAssetFromBucketCommand{BucketID: *idVO, AssetID: input.AssetID}, nil
}
