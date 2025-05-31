# VPC Module Variables
variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
}

variable "private_subnets" {
  description = "List of CIDR blocks for private subnets"
  type        = list(string)
}

#API Module Variables
variable "api_name" {
  description = "Name of the API Gateway"
  type        = string
}

variable "api_key_name" {
  description = "Name of the API Key"
  type        = string
}

variable "api_key_description" {
  description = "Description of the API Key"
  type        = string
}

variable "api_description" {
  description = "Description of the API"
  type        = string
}

variable "api_usage_plan_name" {
  description = "Name of the usage plan"
  type        = string
}

variable "api_usage_plan_description" {
  description = "Description of the usage plan"
  type        = string
}

variable "api_throttling_burst_limit" {
  description = "Burst limit for throttling"
  type        = number
}

variable "api_throttling_rate_limit" {
  description = "Rate limit for throttling (requests per second)"
  type        = number
}

variable "api_quota_limit" {
  description = "Maximum number of requests allowed"
  type        = number
}

variable "api_quota_offset" {
  description = "Offset for the quota limit"
  type        = number
}

variable "api_quota_period" {
  description = "Time period for quota reset (e.g., DAY, WEEK, MONTH)"
  type        = string
}

variable "api_stage_name" {
  description = "Stage name for the deployment"
  type        = string
}

#Storage Module Variables
variable "raw_storage_s3_bucket_name" {
  description = "Name of the S3 bucket for raw video storage"
  type        = string
}

variable "transcoded_storage_s3_bucket_name" {
  description = "Name of the S3 bucket for transcoded video storage"
  type        = string
}

variable "thumbnail_storage_s3_bucket_name" {
  description = "Name of the S3 bucket for thumbnail video storage"
  type        = string
}

#Transcoder Module Variables
variable "sqs_queue_name" {
  description = "The name of the SQS queue for transcoding jobs"
  type        = string
  default     = "TranscodingQueue"
}

variable "sqs_queue_visibility_timeout" {
  description = "The visibility timeout for the SQS queue (in seconds)"
  type        = number
  default     = 300
}

#Common Variables
variable "tags" {
  description = "Tags for resources"
  type        = map(string)
}

variable "aws_region" {
  description = "Name of the AWS region"
  type        = string
}

