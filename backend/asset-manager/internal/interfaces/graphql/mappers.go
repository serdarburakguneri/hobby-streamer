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
	return assetCommands.CreateAssetCommand{
		Slug:      *slug,
		Title:     title,
		AssetType: at,
		OwnerID:   owner,
		ParentID:  parent,
	}, nil
}

func MapPatchAssetInput(id string, patches []*JSONPatch) (assetCommands.PatchAssetCommand, error) {
	idVO, err := assetvo.NewAssetID(id)
	if err != nil {
		return assetCommands.PatchAssetCommand{}, err
	}
	ops := make([]assetCommands.JSONPatchOperation, len(patches))
	for i, p := range patches {
		ops[i] = assetCommands.JSONPatchOperation{Op: p.Op, Path: p.Path, Value: p.Value}
	}
	return assetCommands.PatchAssetCommand{ID: *idVO, Patches: ops}, nil
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
	return bucketCommands.CreateBucketCommand{Name: *input.Name, Key: *input.Key, OwnerID: owner}, nil
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
