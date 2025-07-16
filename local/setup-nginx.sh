#!/bin/bash
set -e

cd "$(dirname "$0")"

echo "[INFO] Starting nginx proxy service..."
docker-compose up -d nginx

echo "[INFO] Nginx proxy is running!"
echo "[INFO] - Nginx (HLS/DASH Proxy): http://localhost:8083" 