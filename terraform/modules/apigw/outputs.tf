output "api_id" {
  description = "ID of the REST API"
  value       = aws_api_gateway_rest_api.api.id
}

output "api_key_id" {
  description = "ID of the API Key"
  value       = aws_api_gateway_api_key.api_key.id
}

output "api_key_value" {
  description = "Value of the API Key"
  value       = aws_api_gateway_api_key.api_key.value
}

output "usage_plan_id" {
  description = "ID of the usage plan"
  value       = aws_api_gateway_usage_plan.usage_plan.id
}
