resource "aws_api_gateway_rest_api" "api" {
  name        = var.api_name
  description = var.api_description
  tags        = var.tags
}

resource "aws_api_gateway_api_key" "api_key" {
  name        = var.api_key_name
  description = var.api_key_description
  enabled     = true
  tags        = var.tags
}

resource "aws_api_gateway_usage_plan" "usage_plan" {
  name        = var.usage_plan_name
  description = var.usage_plan_description

  throttle_settings {
    burst_limit = var.throttling_burst_limit
    rate_limit  = var.throttling_rate_limit
  }

  quota_settings {
    limit  = var.quota_limit
    offset = var.quota_offset
    period = var.quota_period
  }

  api_stages {
    api_id    = aws_api_gateway_rest_api.api.id
    stage     = var.stage_name
  }

  tags = var.tags
}

resource "aws_api_gateway_usage_plan_key" "usage_plan_key" {
  key_id        = aws_api_gateway_api_key.api_key.id
  key_type      = "API_KEY"
  usage_plan_id = aws_api_gateway_usage_plan.usage_plan.id
}