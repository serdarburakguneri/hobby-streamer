#!/bin/bash
set -e

# Ensure we're in the local directory
cd "$(dirname "$0")"

source "setup-environment.sh"

echo "[INFO] Setting up Kafka topics for Hobby Streamer..."

# Wait for Kafka to be ready
echo "[INFO] Waiting for Kafka to be ready..."
until docker exec kafka kafka-topics --bootstrap-server localhost:29092 --list > /dev/null 2>&1; do
    echo "[INFO] Waiting for Kafka to be ready..."
    sleep 5
done

# Create Kafka topics
TOPICS=(
    "asset-events"
    "bucket-events"
    "job-requests"
    "job-completions"
    "content-analysis"
    "raw-video-uploaded"
)

for topic in "${TOPICS[@]}"; do
    echo "[INFO] Creating Kafka topic: $topic"
    docker exec kafka kafka-topics \
        --create \
        --if-not-exists \
        --bootstrap-server localhost:29092 \
        --replication-factor 1 \
        --partitions 3 \
        --topic "$topic"
done

echo "[INFO] Kafka topics setup completed successfully!" 