
resource "aws_s3_bucket_notification" "trigger_transcoder" {
  bucket = var.raw_storage_bucket_id

  queue {
    queue_arn     = var.transcoding_queue_arn
    events        = ["s3:ObjectCreated:*"]
    filter_suffix = ".mp4"
  }
}