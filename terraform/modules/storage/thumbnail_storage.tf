//todo: add policies as other modules start interaction with this s3 bucket

resource "aws_s3_bucket" "thumbnail_storage_bucket" {
  bucket = var.thumbnail_storage_s3_bucket_name
  tags   = var.tags
}

resource "aws_s3_bucket_versioning" "thumbnail_bucket_versioning" {
  bucket = aws_s3_bucket.thumbnail_storage_bucket.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "thumbnail_expiration" {
  bucket = aws_s3_bucket.thumbnail_storage_bucket.id

  rule {
    id     = "expire-old-versions"
    status = "Enabled"

    noncurrent_version_expiration {
      noncurrent_days = 30
    }
  }
}