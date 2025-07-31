#!/bin/bash
set -e

# Ensure we're in the local directory
cd "$(dirname "$0")"

source "setup-environment.sh"

echo "[INFO] Testing Kafka topics for Hobby Streamer..."

# Wait for Kafka to be ready
echo "[INFO] Waiting for Kafka to be ready..."
until docker exec kafka kafka-topics --bootstrap-server localhost:29092 --list > /dev/null 2>&1; do
    echo "[INFO] Waiting for Kafka to be ready..."
    sleep 5
done

# Test Kafka topics
TOPICS=(
    "asset-events"
    "bucket-events"
    "job-requests"
    "job-completions"
    "content-analysis"
    "raw-video-uploaded"
)

for topic in "${TOPICS[@]}"; do
    echo "[INFO] Testing Kafka topic: $topic"
    if docker exec kafka kafka-topics --bootstrap-server localhost:29092 --describe --topic "$topic" > /dev/null 2>&1; then
        echo "[INFO] Topic $topic exists and is accessible"
    else
        echo "[ERROR] Topic $topic is not accessible"
        exit 1
    fi
done

echo "[INFO] Kafka topic testing completed successfully!" 