output "lambda_arn" {
  description = "ARN of the Lambda function"
  value       = aws_lambda_function.generate_url.arn
}

output "raw_storage_bucket_id" {
  value = aws_s3_bucket.raw_storage_bucket.id
}

output "raw_storage_bucket_arn" {
  value = aws_s3_bucket.raw_storage_bucket.arn
}