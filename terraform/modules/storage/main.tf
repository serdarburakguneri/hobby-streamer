resource "aws_s3_bucket" "video_bucket" {
  bucket = var.s3_bucket_name

  tags = var.tags
}

resource "aws_s3_bucket_versioning" "video_bucket_versioning" {
  bucket = aws_s3_bucket.video_bucket.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_policy" "bucket_policy" {
  bucket = aws_s3_bucket.video_bucket.id

  policy = jsonencode({
    Version   = "2012-10-17",
    Statement = [
      {
        Sid      = "AllowLambdaAccess",
        Effect   = "Allow",
        Principal = { AWS = aws_iam_role.lambda_role.arn },
        Action   = ["s3:PutObject", "s3:GetObject"],
        Resource = [
          "${aws_s3_bucket.video_bucket.arn}",
          "${aws_s3_bucket.video_bucket.arn}/*"
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
        Effect = "Allow",
        Principal = { Service = "lambda.amazonaws.com" },
        Action = "sts:AssumeRole"
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
        Action   = ["s3:PutObject", "s3:GetObject"],
        Resource = [
          "${aws_s3_bucket.video_bucket.arn}",
          "${aws_s3_bucket.video_bucket.arn}/*"
        ]
      },
      {
        Effect   = "Allow",
        Action   = ["logs:CreateLogGroup", "logs:CreateLogStream", "logs:PutLogEvents"],
        Resource = "arn:aws:logs:*:*:*"
      }
    ]
  })
}

resource "aws_lambda_function" "generate_url" {
  function_name = "generate-presigned-url"
  role          = aws_iam_role.lambda_role.arn
  runtime = "provided.al2"  # Amazon Linux 2 runtime
  handler = "bootstrap"    # Name of the binary (required for Go)
  filename      = "${path.module}/lambda/main.zip"
  timeout       = 10
  memory_size   = 128

  environment {
    variables = {
      BUCKET_NAME = aws_s3_bucket.video_bucket.id,
      BUCKET_REGION  = var.aws_region
    }
  }

  tags = var.tags
}


