output "s3_bucket_id" {
  description = "ID of the S3 bucket"
  value       = aws_s3_bucket.video_bucket.id
}

output "lambda_arn" {
  description = "ARN of the Lambda function"
  value       = aws_lambda_function.generate_url.arn
}
