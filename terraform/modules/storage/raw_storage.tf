resource "aws_s3_bucket" "raw_storage_bucket" {
  bucket = var.raw_storage_s3_bucket_name
  tags   = var.tags
}

resource "aws_s3_bucket_versioning" "video_bucket_versioning" {
  bucket = aws_s3_bucket.raw_storage_bucket.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "raw_lifecycle" {
  bucket = aws_s3_bucket.raw_storage_bucket.id

  rule {
    id     = "expire-temp-uploads"
    status = "Enabled"

    expiration {
      days = 7
    }

    filter {} # Applies to all objects in the bucket
  }
}

resource "aws_s3_bucket_policy" "bucket_policy" {
  bucket = aws_s3_bucket.raw_storage_bucket.id

  policy = jsonencode({
    Version   = "2012-10-17",
    Statement = [
      {
        Sid       = "AllowLambdaAccess",
        Effect    = "Allow",
        Principal = {
          AWS = aws_iam_role.lambda_role.arn
        },
        Action = [
          "s3:PutObject",
          "s3:GetObject"
        ],
        Resource = [
          "${aws_s3_bucket.raw_storage_bucket.arn}",
          "${aws_s3_bucket.raw_storage_bucket.arn}/*"
        ]
      }
    ]
  })
}

resource "aws_iam_role" "lambda_role" {
  name = "lambda-storage-role"

  assume_role_policy = jsonencode({
    Version   = "2012-10-17",
    Statement = [
      {
        Effect    = "Allow",
        Principal = { Service = "lambda.amazonaws.com" },
        Action    = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy" "lambda_policy" {
  role = aws_iam_role.lambda_role.id

  policy = jsonencode({
    Version   = "2012-10-17",
    Statement = [
      {
        Effect   = "Allow",
        Action   = [
          "s3:PutObject",
          "s3:GetObject"
        ],
        Resource = [
          "${aws_s3_bucket.raw_storage_bucket.arn}",
          "${aws_s3_bucket.raw_storage_bucket.arn}/*"
        ]
      },
      {
        Effect   = "Allow",
        Action   = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        Resource = "arn:aws:logs:*:*:*"
      }
    ]
  })
}

resource "aws_lambda_function" "generate_url" {
  function_name    = "generate-presigned-url"
  role             = aws_iam_role.lambda_role.arn
  runtime          = "provided.al2"
  handler          = "bootstrap"
  filename         = "${path.module}/../build/generate_presigned_upload_url.zip"
  source_code_hash = filebase64sha256("${path.module}/../build/generate_presigned_upload_url.zip")
  timeout          = 10
  memory_size      = 128

  environment {
    variables = {
      BUCKET_NAME   = aws_s3_bucket.raw_storage_bucket.bucket
      BUCKET_REGION = var.aws_region
    }
  }

  tags = var.tags
}

resource "aws_api_gateway_resource" "generate_url" {
  rest_api_id = var.api_id
  parent_id   = var.api_root_resource_id
  path_part   = "generate-presigned-upload-url"
}

resource "aws_api_gateway_method" "generate_url_method" {
  rest_api_id      = var.api_id
  resource_id      = aws_api_gateway_resource.generate_url.id
  http_method      = "POST"
  authorization    = "NONE"
  api_key_required = true
}

resource "aws_api_gateway_integration" "lambda_integration" {
  rest_api_id             = var.api_id
  resource_id             = aws_api_gateway_resource.generate_url.id
  integration_http_method = "POST"
  http_method             = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.generate_url.invoke_arn

  depends_on = [aws_api_gateway_method.generate_url_method]
}

resource "aws_lambda_permission" "apigw_permission" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.generate_url.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "arn:aws:execute-api:${var.aws_region}:${var.account_id}:${var.api_id}/${var.stage_name}/POST/generate-presigned-upload-url"
}