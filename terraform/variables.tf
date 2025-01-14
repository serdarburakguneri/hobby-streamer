# VPC Module Variables
variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
}

variable "private_subnets" {
  description = "List of CIDR blocks for private subnets"
  type        = list(string)
}

# APIGW Module Variables
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

variable "usage_plan_name" {
  description = "Name of the usage plan"
  type        = string
}

variable "usage_plan_description" {
  description = "Description of the usage plan"
  type        = string
}

variable "throttling_burst_limit" {
  description = "Burst limit for throttling"
  type        = number
}

variable "throttling_rate_limit" {
  description = "Rate limit for throttling (requests per second)"
  type        = number
}

variable "quota_limit" {
  description = "Maximum number of requests allowed"
  type        = number
}

variable "quota_offset" {
  description = "Offset for the quota limit"
  type        = number
}

variable "quota_period" {
  description = "Time period for quota reset (e.g., DAY, WEEK, MONTH)"
  type        = string
}

#Storage Module Variables
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

