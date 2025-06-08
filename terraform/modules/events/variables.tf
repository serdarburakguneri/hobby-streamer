variable "raw_storage_bucket_id" {
  description = "The ID of the raw S3 bucket (for bucket notifications)"
  type        = string
}

variable "transcoding_queue_arn" {
  description = "The ARN of the transcoder SQS queue"
  type        = string
}