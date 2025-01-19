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

variable "sqs_queue_name" {
  description = "The name of the SQS queue for transcoding jobs"
  type        = string
  default     = "TranscodingQueue"
}

variable "sqs_queue_visibility_timeout" {
  description = "The visibility timeout for the SQS queue (in seconds)"
  type        = number
  default     = 300  # Adjust based on average transcoding time
}

variable "tags" {
  description = "Tags for resources"
  type        = map(string)
}
