variable "s3_bucket_name" {
  description = "Name of the S3 bucket for video storage"
  type        = string
}

variable "tags" {
  description = "Tags for resources"
  type        = map(string)
}

variable "aws_region" {
  description = "Name of the AWS region"
  type        = string
}