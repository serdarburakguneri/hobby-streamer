#!/bin/bash
set -e

cd "$(dirname "$0")"

source "setup-environment.sh"

echo "[INFO] Setting up CloudFront distributions for streaming buckets..."

if ! aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT cloudfront list-distributions --region $AWS_REGION > /dev/null 2>&1; then
    echo "[INFO] CloudFront not available in LocalStack, creating mock distributions..."
    
    echo "[INFO] Creating CloudFront distribution for HLS streaming..."
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT cloudfront create-distribution \
        --distribution-config '{
            "CallerReference": "hls-streaming-'$(date +%s)'",
            "Comment": "HLS Streaming Distribution",
            "DefaultCacheBehavior": {
                "TargetOriginId": "hls-storage-origin",
                "ViewerProtocolPolicy": "allow-all",
                "TrustedSigners": {
                    "Enabled": false,
                    "Quantity": 0
                },
                "ForwardedValues": {
                    "QueryString": false,
                    "Cookies": {
                        "Forward": "none"
                    }
                },
                "MinTTL": 0,
                "DefaultTTL": 86400,
                "MaxTTL": 31536000
            },
            "Enabled": true,
            "Origins": {
                "Quantity": 1,
                "Items": [
                    {
                        "Id": "hls-storage-origin",
                        "DomainName": "hls-storage.s3.localhost.localstack.cloud",
                        "S3OriginConfig": {
                            "OriginAccessIdentity": ""
                        }
                    }
                ]
            },
            "Aliases": {
                "Quantity": 1,
                "Items": ["hls-streaming.localhost"]
            }
        }' \
        --region $AWS_REGION > /dev/null 2>&1 || echo "[WARN] CloudFront not supported in LocalStack"

    echo "[INFO] Creating CloudFront distribution for DASH streaming..."
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT cloudfront create-distribution \
        --distribution-config '{
            "CallerReference": "dash-streaming-'$(date +%s)'",
            "Comment": "DASH Streaming Distribution",
            "DefaultCacheBehavior": {
                "TargetOriginId": "dash-storage-origin",
                "ViewerProtocolPolicy": "allow-all",
                "TrustedSigners": {
                    "Enabled": false,
                    "Quantity": 0
                },
                "ForwardedValues": {
                    "QueryString": false,
                    "Cookies": {
                        "Forward": "none"
                    }
                },
                "MinTTL": 0,
                "DefaultTTL": 86400,
                "MaxTTL": 31536000
            },
            "Enabled": true,
            "Origins": {
                "Quantity": 1,
                "Items": [
                    {
                        "Id": "dash-storage-origin",
                        "DomainName": "dash-storage.s3.localhost.localstack.cloud",
                        "S3OriginConfig": {
                            "OriginAccessIdentity": ""
                        }
                    }
                ]
            },
            "Aliases": {
                "Quantity": 1,
                "Items": ["dash-streaming.localhost"]
            }
        }' \
        --region $AWS_REGION > /dev/null 2>&1 || echo "[WARN] CloudFront not supported in LocalStack"

    echo "[INFO] Creating CloudFront distribution for images..."
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT cloudfront create-distribution \
        --distribution-config '{
            "CallerReference": "images-'$(date +%s)'",
            "Comment": "Images Distribution",
            "DefaultCacheBehavior": {
                "TargetOriginId": "images-storage-origin",
                "ViewerProtocolPolicy": "allow-all",
                "TrustedSigners": {
                    "Enabled": false,
                    "Quantity": 0
                },
                "ForwardedValues": {
                    "QueryString": false,
                    "Cookies": {
                        "Forward": "none"
                    }
                },
                "MinTTL": 0,
                "DefaultTTL": 86400,
                "MaxTTL": 31536000
            },
            "Enabled": true,
            "Origins": {
                "Quantity": 1,
                "Items": [
                    {
                        "Id": "images-storage-origin",
                        "DomainName": "images-storage.s3.localhost.localstack.cloud",
                        "S3OriginConfig": {
                            "OriginAccessIdentity": ""
                        }
                    }
                ]
            },
            "Aliases": {
                "Quantity": 1,
                "Items": ["images.localhost"]
            }
        }' \
        --region $AWS_REGION > /dev/null 2>&1 || echo "[WARN] CloudFront not supported in LocalStack"
else
    echo "[INFO] CloudFront is available, creating real distributions..."
    
    echo "[INFO] Creating CloudFront distribution for HLS streaming..."
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT cloudfront create-distribution \
        --distribution-config '{
            "CallerReference": "hls-streaming-'$(date +%s)'",
            "Comment": "HLS Streaming Distribution",
            "DefaultCacheBehavior": {
                "TargetOriginId": "hls-storage-origin",
                "ViewerProtocolPolicy": "allow-all",
                "TrustedSigners": {
                    "Enabled": false,
                    "Quantity": 0
                },
                "ForwardedValues": {
                    "QueryString": false,
                    "Cookies": {
                        "Forward": "none"
                    }
                },
                "MinTTL": 0,
                "DefaultTTL": 86400,
                "MaxTTL": 31536000
            },
            "Enabled": true,
            "Origins": {
                "Quantity": 1,
                "Items": [
                    {
                        "Id": "hls-storage-origin",
                        "DomainName": "hls-storage.s3.'$AWS_REGION'.amazonaws.com",
                        "S3OriginConfig": {
                            "OriginAccessIdentity": ""
                        }
                    }
                ]
            }
        }' \
        --region $AWS_REGION

    echo "[INFO] Creating CloudFront distribution for DASH streaming..."
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT cloudfront create-distribution \
        --distribution-config '{
            "CallerReference": "dash-streaming-'$(date +%s)'",
            "Comment": "DASH Streaming Distribution",
            "DefaultCacheBehavior": {
                "TargetOriginId": "dash-storage-origin",
                "ViewerProtocolPolicy": "allow-all",
                "TrustedSigners": {
                    "Enabled": false,
                    "Quantity": 0
                },
                "ForwardedValues": {
                    "QueryString": false,
                    "Cookies": {
                        "Forward": "none"
                    }
                },
                "MinTTL": 0,
                "DefaultTTL": 86400,
                "MaxTTL": 31536000
            },
            "Enabled": true,
            "Origins": {
                "Quantity": 1,
                "Items": [
                    {
                        "Id": "dash-storage-origin",
                        "DomainName": "dash-storage.s3.'$AWS_REGION'.amazonaws.com",
                        "S3OriginConfig": {
                            "OriginAccessIdentity": ""
                        }
                    }
                ]
            }
        }' \
        --region $AWS_REGION

    echo "[INFO] Creating CloudFront distribution for images..."
    aws --endpoint-url=$LOCALSTACK_EXTERNAL_ENDPOINT cloudfront create-distribution \
        --distribution-config '{
            "CallerReference": "images-'$(date +%s)'",
            "Comment": "Images Distribution",
            "DefaultCacheBehavior": {
                "TargetOriginId": "images-storage-origin",
                "ViewerProtocolPolicy": "allow-all",
                "TrustedSigners": {
                    "Enabled": false,
                    "Quantity": 0
                },
                "ForwardedValues": {
                    "QueryString": false,
                    "Cookies": {
                        "Forward": "none"
                    }
                },
                "MinTTL": 0,
                "DefaultTTL": 86400,
                "MaxTTL": 31536000
            },
            "Enabled": true,
            "Origins": {
                "Quantity": 1,
                "Items": [
                    {
                        "Id": "images-storage-origin",
                        "DomainName": "images-storage.s3.'$AWS_REGION'.amazonaws.com",
                        "S3OriginConfig": {
                            "OriginAccessIdentity": ""
                        }
                    }
                ]
            }
        }' \
        --region $AWS_REGION
fi

echo "[INFO] CloudFront distributions setup completed!"
echo "[INFO] Note: For production, you'll need to configure real CloudFront distributions in AWS" 