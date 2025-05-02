
resource "aws_sqs_queue" "dlq" {
  name                        = "${var.sqs_queue_name}_DLQ"
  visibility_timeout_seconds  = var.sqs_queue_visibility_timeout
  message_retention_seconds   = 1209600
  tags                        = var.tags
}

resource "aws_sqs_queue" "transcoding_queue" {
  name                        = var.sqs_queue_name
  visibility_timeout_seconds  = var.sqs_queue_visibility_timeout
  message_retention_seconds   = 1209600

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.dlq.arn,
    maxReceiveCount     = 5
  })

  tags = var.tags
}

resource "aws_iam_role" "apigw_sqs_role" {
  name = "apigw-sqs-integration-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Service = "apigateway.amazonaws.com"
        },
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy" "sqs_policy" {
  role = aws_iam_role.apigw_sqs_role.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = ["sqs:SendMessage"],
        Resource = aws_sqs_queue.transcoding_queue.arn
      }
    ]
  })
}

resource "aws_api_gateway_resource" "transcode_resource" {
  rest_api_id = var.api_id
  parent_id   = var.api_root_resource_id
  path_part   = "transcode"
}

resource "aws_api_gateway_method" "post_method" {
  rest_api_id      = var.api_id
  resource_id      = aws_api_gateway_resource.transcode_resource.id
  http_method      = "POST"
  authorization    = "NONE"
  api_key_required = true
}

resource "aws_api_gateway_integration" "post_sqs_integration" {
  rest_api_id             = var.api_id
  resource_id             = aws_api_gateway_resource.transcode_resource.id
  http_method             = "POST"
  integration_http_method = "POST"
  type                    = "AWS"
  uri                     = "arn:aws:apigateway:${var.aws_region}:sqs:path/${aws_sqs_queue.transcoding_queue.name}"
  credentials             = aws_iam_role.apigw_sqs_role.arn

  request_templates = {
    "application/json" = <<EOF
    {
      "QueueUrl": "https://sqs.${var.aws_region}.amazonaws.com/${var.account_id}/${aws_sqs_queue.transcoding_queue.name}",
      "MessageBody": "$util.escapeJavaScript($input.body)"
    }
    EOF
  }

  depends_on = [
    aws_api_gateway_method.post_method,
    aws_iam_role_policy.sqs_policy
  ]
}