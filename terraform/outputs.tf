output "vpc_id" {
  description = "The ID of the VPC"
  value       = module.vpc.vpc_id
}

output "private_subnet_ids" {
  description = "List of private subnet IDs"
  value       = module.vpc.private_subnet_ids
}

output "storage_api_id" {
  description = "ID of the storage REST API"
  value       = module.storage-api.api_id
}

output "api_key_value" {
  description = "The value of the API Key"
  value       = module.storage-api.api_key_value
  sensitive   = true
}