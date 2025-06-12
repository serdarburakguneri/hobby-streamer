#!/bin/bash
set -e

BUILD_DIR=./terraform/modules/build
mkdir -p "$BUILD_DIR"

# Initialize go.mod if not present
if [ ! -f "go.mod" ]; then
  echo "⚙️  No go.mod found, initializing Go module..."
  go mod init github.com/serdarburakguneri/hobby-streamer
  go mod tidy
fi

# Format: zip_name relative_cmd_path_from_root
BUILD_TARGETS=(
  "generate_presigned_upload_url services/storage/cmd/generate_presigned_upload_url"
  "save_asset services/asset-manager/cmd/save_asset"
  "get_asset services/asset-manager/cmd/get_asset"
  "list_assets services/asset-manager/cmd/list_assets"
)

for entry in "${BUILD_TARGETS[@]}"; do
  ZIP_NAME=$(echo "$entry" | awk '{print $1}')
  CMD_PATH=$(echo "$entry" | awk '{print $2}')
  OUTPUT_DIR="$BUILD_DIR/$ZIP_NAME"
  mkdir -p "$OUTPUT_DIR"

  echo "📦 Building $ZIP_NAME from $CMD_PATH..."

  GOOS=linux GOARCH=arm64 go build -tags lambda.norpc \
    -o "$OUTPUT_DIR/bootstrap" \
    "./$CMD_PATH"

  (cd "$OUTPUT_DIR" && zip "../${ZIP_NAME}.zip" bootstrap)
done