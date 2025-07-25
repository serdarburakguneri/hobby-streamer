#!/bin/bash
set -e

# Ensure we're in the local directory
cd "$(dirname "$0")"

source "setup-environment.sh"

echo "[INFO] Setting up SQS queues for Hobby Streamer..."

for queue in job-queue-dlq completion-queue-dlq asset-events-dlq; do
    if ! aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-url --queue-name $queue --region $AWS_REGION > /dev/null 2>&1; then
        echo "[INFO] Creating SQS DLQ: $queue"
        aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs create-queue --queue-name $queue --region $AWS_REGION > /dev/null
    else
        echo "[INFO] SQS DLQ $queue already exists."
    fi
done

JOB_QUEUE_DLQ_ARN="arn:aws:sqs:us-east-1:000000000000:job-queue-dlq"
COMPLETION_QUEUE_DLQ_ARN="arn:aws:sqs:us-east-1:000000000000:completion-queue-dlq"
ASSET_EVENTS_DLQ_ARN="arn:aws:sqs:us-east-1:000000000000:asset-events-dlq"

if ! aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-url --queue-name job-queue --region $AWS_REGION > /dev/null 2>&1; then
    echo "[INFO] Creating SQS queue: job-queue"
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs create-queue --queue-name job-queue --region $AWS_REGION > /dev/null
    echo "[INFO] Setting redrive policy for job-queue"
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs set-queue-attributes \
        --queue-url http://localhost:4566/000000000000/job-queue \
        --attributes '{"RedrivePolicy":"{\"deadLetterTargetArn\":\"'$JOB_QUEUE_DLQ_ARN'\",\"maxReceiveCount\":\"3\"}"}' \
        --region $AWS_REGION > /dev/null
else
    echo "[INFO] SQS queue job-queue already exists."
fi

if ! aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-url --queue-name completion-queue --region $AWS_REGION > /dev/null 2>&1; then
    echo "[INFO] Creating SQS queue: completion-queue"
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs create-queue --queue-name completion-queue --region $AWS_REGION > /dev/null
    echo "[INFO] Setting redrive policy for completion-queue"
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs set-queue-attributes \
        --queue-url http://localhost:4566/000000000000/completion-queue \
        --attributes '{"RedrivePolicy":"{\"deadLetterTargetArn\":\"'$COMPLETION_QUEUE_DLQ_ARN'\",\"maxReceiveCount\":\"3\"}"}' \
        --region $AWS_REGION > /dev/null
else
    echo "[INFO] SQS queue completion-queue already exists."
fi

if ! aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs get-queue-url --queue-name asset-events --region $AWS_REGION > /dev/null 2>&1; then
    echo "[INFO] Creating SQS queue: asset-events"
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs create-queue --queue-name asset-events --region $AWS_REGION > /dev/null
    echo "[INFO] Setting redrive policy for asset-events"
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT sqs set-queue-attributes \
        --queue-url http://localhost:4566/000000000000/asset-events \
        --attributes '{"RedrivePolicy":"{\"deadLetterTargetArn\":\"'$ASSET_EVENTS_DLQ_ARN'\",\"maxReceiveCount\":\"3\"}"}' \
        --region $AWS_REGION > /dev/null
else
    echo "[INFO] SQS queue asset-events already exists."
fi

echo "[INFO] SQS queues setup completed successfully!" 