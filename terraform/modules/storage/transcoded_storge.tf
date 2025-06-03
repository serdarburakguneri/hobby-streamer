//todo: add policies as other modules start interaction with this s3 bucket

resource "aws_s3_bucket" "transcoded_storage_bucket" {
  bucket = var.transcoded_storage_s3_bucket_name
  tags   = var.tags
}

resource "aws_s3_bucket_versioning" "transcoded_bucket_versioning" {
  bucket = aws_s3_bucket.transcoded_storage_bucket.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "transcoded_expiration" {
  bucket = aws_s3_bucket.transcoded_storage_bucket.id

  rule {
    id     = "expire-old-versions"
    status = "Enabled"

    noncurrent_version_expiration {
      days = 30
    }
  }
}