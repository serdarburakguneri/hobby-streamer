#!/bin/bash
set -e

BUILD_DIR=../../build
mkdir -p $BUILD_DIR

for cmd in generate_upload_url; do
  echo "Building $cmd..."
  GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o $BUILD_DIR/${cmd}/bootstrap ./cmd/${cmd}
  (cd $BUILD_DIR/${cmd} && zip ../${cmd}.zip bootstrap)
done