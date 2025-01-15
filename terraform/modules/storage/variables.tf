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

variable "api_id" {
  description = "Api id"
  type        = string
}

variable "api_root_resource_id" {
  description = "Api root resource id"
  type        = string
}

variable "stage_name" {
  description = "deployment stage name"
  type        = string
}

variable "account_id" {
  description = "AWS Account ID"
  type        = string
}

