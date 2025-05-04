output "transcoding_queue_url" {
  description = "The URL of the transcoder SQS queue"
  value       = aws_sqs_queue.transcoding_queue.id
}

output "transcoding_queue_arn" {
  description = "The ARN of the transcoder SQS queue"
  value       = aws_sqs_queue.transcoding_queue.arn
}