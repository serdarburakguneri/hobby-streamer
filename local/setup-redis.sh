#!/bin/bash
set -e

cd "$(dirname "$0")"

echo "[INFO] Starting Redis..."
docker-compose up -d redis

echo "[INFO] Redis is running at redis://localhost:6379" 