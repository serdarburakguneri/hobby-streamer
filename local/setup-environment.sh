#!/bin/bash

set -e

echo "[INFO] Setting up environment variables for Hobby Streamer..."

export AWS_REGION="us-east-1"
export AWS_ACCESS_KEY_ID="test"
export AWS_SECRET_ACCESS_KEY="test"
export LOCALSTACK_EXTERNAL_ENDPOINT="http://localhost:4566"
export LOCALSTACK_INTERNAL_ENDPOINT="http://localstack:4566"

export JOB_QUEUE_URL="http://localstack:4566/000000000000/job-queue"
export COMPLETION_QUEUE_URL="http://localstack:4566/000000000000/completion-queue"

export NEO4J_URI="bolt://neo4j:7687"
export NEO4J_USERNAME="neo4j"
export NEO4J_PASSWORD="password"

export KEYCLOAK_URL="https://keycloak:8443"
export KEYCLOAK_REALM="hobby"
export KEYCLOAK_CLIENT_ID="asset-manager"
export KEYCLOAK_CLIENT_SECRET="your-client-secret"

export REDIS_URL="redis://redis:6379"

export CONTENT_BUCKET="content-east"
export CDN_PREFIX="http://localhost:8083/cdn"

# Kafka Configuration
export KAFKA_BOOTSTRAP_SERVERS="localhost:9092"
export KAFKA_INTERNAL_BOOTSTRAP_SERVERS="kafka:29092"
export KAFKA_TOPIC_ASSET_EVENTS="asset-events"
export KAFKA_TOPIC_BUCKET_EVENTS="bucket-events"
export KAFKA_TOPIC_JOB_REQUESTS="job-requests"
export KAFKA_TOPIC_JOB_COMPLETIONS="job-completions"
export KAFKA_TOPIC_CONTENT_ANALYSIS="content-analysis"

echo "[INFO] Environment variables set successfully!"
echo "[INFO] AWS Region: $AWS_REGION"
echo "[INFO] LocalStack External Endpoint: $LOCALSTACK_EXTERNAL_ENDPOINT"
echo "[INFO] Job Queue URL: $JOB_QUEUE_URL"
echo "[INFO] Completion Queue URL: $COMPLETION_QUEUE_URL"
echo "[INFO] Neo4j URI: $NEO4J_URI"
echo "[INFO] Keycloak URL: $KEYCLOAK_URL"
echo "[INFO] Redis URL: $REDIS_URL"
echo "[INFO] Content Bucket: $CONTENT_BUCKET"
echo "[INFO] CDN Prefix: $CDN_PREFIX"
echo "[INFO] Kafka Bootstrap Servers: $KAFKA_BOOTSTRAP_SERVERS"
echo "[INFO] Kafka Internal Bootstrap Servers: $KAFKA_INTERNAL_BOOTSTRAP_SERVERS" 