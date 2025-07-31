#!/bin/bash
set -e

cd "$(dirname "$0")"

source "setup-environment.sh"

echo "[INFO] Setting up Kafka infrastructure for Hobby Streamer..."

echo "[INFO] Waiting for Zookeeper to be ready..."
max_attempts=60
attempt=0
while [ $attempt -lt $max_attempts ]; do
  if docker exec zookeeper echo ruok | nc localhost 2181 > /dev/null 2>&1; then
    echo "[INFO] Zookeeper is ready."
    break
  fi
  echo "[INFO] Zookeeper not ready yet, waiting... (attempt $((attempt + 1))/$max_attempts)"
  sleep 5
  attempt=$((attempt + 1))
done

if [ $attempt -eq $max_attempts ]; then
  echo "[ERROR] Zookeeper failed to become ready after $max_attempts attempts"
  exit 1
fi

echo "[INFO] Waiting for Kafka to be ready..."
max_attempts=60
attempt=0
while [ $attempt -lt $max_attempts ]; do
  if docker exec kafka kafka-topics --bootstrap-server localhost:9092 --list > /dev/null 2>&1; then
    echo "[INFO] Kafka is ready."
    break
  fi
  echo "[INFO] Kafka not ready yet, waiting... (attempt $((attempt + 1))/$max_attempts)"
  sleep 5
  attempt=$((attempt + 1))
done

if [ $attempt -eq $max_attempts ]; then
  echo "[ERROR] Kafka failed to become ready after $max_attempts attempts"
  exit 1
fi

echo "[INFO] Waiting additional time for Kafka to fully stabilize..."
sleep 10

echo "[INFO] Creating Kafka topics..."

# Asset Events Topic (partitioned by assetId for ordering)
echo "[INFO] Creating asset-events topic..."
docker exec kafka kafka-topics \
  --bootstrap-server localhost:9092 \
  --create \
  --topic asset-events \
  --partitions 6 \
  --replication-factor 1 \
  --config retention.ms=604800000 \
  --config cleanup.policy=delete \
  --config compression.type=snappy \
  --if-not-exists

# Bucket Events Topic (partitioned by bucketId for ordering)
echo "[INFO] Creating bucket-events topic..."
docker exec kafka kafka-topics \
  --bootstrap-server localhost:9092 \
  --create \
  --topic bucket-events \
  --partitions 3 \
  --replication-factor 1 \
  --config retention.ms=604800000 \
  --config cleanup.policy=delete \
  --config compression.type=snappy \
  --if-not-exists

# Job Request Topics (specific job types)
echo "[INFO] Creating analyze.job.requested topic..."
docker exec kafka kafka-topics \
  --bootstrap-server localhost:9092 \
  --create \
  --topic analyze.job.requested \
  --partitions 4 \
  --replication-factor 1 \
  --config retention.ms=259200000 \
  --config cleanup.policy=delete \
  --config compression.type=snappy \
  --if-not-exists

echo "[INFO] Creating hls.job.requested topic..."
docker exec kafka kafka-topics \
  --bootstrap-server localhost:9092 \
  --create \
  --topic hls.job.requested \
  --partitions 4 \
  --replication-factor 1 \
  --config retention.ms=259200000 \
  --config cleanup.policy=delete \
  --config compression.type=snappy \
  --if-not-exists

# Job Completion Topics (specific job types)
echo "[INFO] Creating analyze.job.completed topic..."
docker exec kafka kafka-topics \
  --bootstrap-server localhost:9092 \
  --create \
  --topic analyze.job.completed \
  --partitions 6 \
  --replication-factor 1 \
  --config retention.ms=604800000 \
  --config cleanup.policy=delete \
  --config compression.type=snappy \
  --if-not-exists

echo "[INFO] Creating hls.job.completed topic..."
docker exec kafka kafka-topics \
  --bootstrap-server localhost:9092 \
  --create \
  --topic hls.job.completed \
  --partitions 6 \
  --replication-factor 1 \
  --config retention.ms=604800000 \
  --config cleanup.policy=delete \
  --config compression.type=snappy \
  --if-not-exists

echo "[INFO] Creating dash.job.completed topic..."
docker exec kafka kafka-topics \
  --bootstrap-server localhost:9092 \
  --create \
  --topic dash.job.completed \
  --partitions 6 \
  --replication-factor 1 \
  --config retention.ms=604800000 \
  --config cleanup.policy=delete \
  --config compression.type=snappy \
  --if-not-exists

# Video Upload Topic
echo "[INFO] Creating raw-video-uploaded topic..."
docker exec kafka kafka-topics \
  --bootstrap-server localhost:9092 \
  --create \
  --topic raw-video-uploaded \
  --partitions 4 \
  --replication-factor 1 \
  --config retention.ms=604800000 \
  --config cleanup.policy=delete \
  --config compression.type=snappy \
  --if-not-exists

# Content Analysis Topics (for future content analyzer service)
echo "[INFO] Creating content-analysis topic..."
docker exec kafka kafka-topics \
  --bootstrap-server localhost:9092 \
  --create \
  --topic content-analysis \
  --partitions 4 \
  --replication-factor 1 \
  --config retention.ms=2592000000 \
  --config cleanup.policy=delete \
  --config compression.type=snappy \
  --if-not-exists

echo "[INFO] Creating content.analysis.requested topic..."
docker exec kafka kafka-topics \
  --bootstrap-server localhost:9092 \
  --create \
  --topic content.analysis.requested \
  --partitions 4 \
  --replication-factor 1 \
  --config retention.ms=259200000 \
  --config cleanup.policy=delete \
  --config compression.type=snappy \
  --if-not-exists

echo "[INFO] Creating content.analysis.completed topic..."
docker exec kafka kafka-topics \
  --bootstrap-server localhost:9092 \
  --create \
  --topic content.analysis.completed \
  --partitions 6 \
  --replication-factor 1 \
  --config retention.ms=604800000 \
  --config cleanup.policy=delete \
  --config compression.type=snappy \
  --if-not-exists

echo "[INFO] Creating content.analysis.failed topic..."
docker exec kafka kafka-topics \
  --bootstrap-server localhost:9092 \
  --create \
  --topic content.analysis.failed \
  --partitions 6 \
  --replication-factor 1 \
  --config retention.ms=604800000 \
  --config cleanup.policy=delete \
  --config compression.type=snappy \
  --if-not-exists

echo "[INFO] Listing all topics..."
docker exec kafka kafka-topics \
  --bootstrap-server localhost:9092 \
  --list

echo "[INFO] Topic configurations:"
for topic in asset-events bucket-events analyze.job.requested hls.job.requested analyze.job.completed hls.job.completed dash.job.completed raw-video-uploaded content-analysis content.analysis.requested content.analysis.completed content.analysis.failed; do
  echo "[INFO] Configuration for $topic:"
  docker exec kafka kafka-topics \
    --bootstrap-server localhost:9092 \
    --describe \
    --topic "$topic"
done

echo "[INFO] Setting up consumer groups for testing..."

# Create test consumer groups
echo "[INFO] Creating test consumer groups..."
docker exec kafka kafka-consumer-groups \
  --bootstrap-server localhost:9092 \
  --group asset-manager-group \
  --describe > /dev/null 2>&1 || echo "[INFO] asset-manager-group will be created when first consumer connects"

docker exec kafka kafka-consumer-groups \
  --bootstrap-server localhost:9092 \
  --group transcoder-group \
  --describe > /dev/null 2>&1 || echo "[INFO] transcoder-group will be created when first consumer connects"

docker exec kafka kafka-consumer-groups \
  --bootstrap-server localhost:9092 \
  --group content-analyzer-group \
  --describe > /dev/null 2>&1 || echo "[INFO] content-analyzer-group will be created when first consumer connects"

echo "[INFO] Kafka infrastructure setup completed successfully!"
echo "[INFO] Kafka UI available at: http://localhost:8086"
echo "[INFO] Kafka broker available at: localhost:9092"
echo ""
echo "[INFO] Topic Summary:"
echo "  - asset-events: 6 partitions, 7 days retention"
echo "  - bucket-events: 3 partitions, 7 days retention"
echo "  - analyze.job.requested: 4 partitions, 3 days retention"
echo "  - hls.job.requested: 4 partitions, 3 days retention"
echo "  - analyze.job.completed: 6 partitions, 7 days retention"
echo "  - hls.job.completed: 6 partitions, 7 days retention"
echo "  - dash.job.completed: 6 partitions, 7 days retention"
echo "  - raw-video-uploaded: 4 partitions, 7 days retention"
echo "  - content-analysis: 4 partitions, 30 days retention"
echo "  - content.analysis.requested: 4 partitions, 3 days retention"
echo "  - content.analysis.completed: 6 partitions, 7 days retention"
echo "  - content.analysis.failed: 6 partitions, 7 days retention"
echo ""
echo "[INFO] Consumer Groups:"
echo "  - asset-manager-group (for asset-manager service)"
echo "  - transcoder-group (for transcoder service)"
echo "  - content-analyzer-group (for future content analyzer)" 