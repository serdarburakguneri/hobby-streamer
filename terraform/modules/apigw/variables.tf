variable "api_name" {
  description = "Name of the REST API"
  type        = string
}

variable "api_key_name" {
  description = "Name of the API Key"
  type        = string
}

variable "api_key_description" {
  description = "Description of the API Key"
  type        = string
}

variable "usage_plan_name" {
  description = "Name of the usage plan"
  type        = string
}

variable "usage_plan_description" {
  description = "Description of the usage plan"
  type        = string
}

variable "throttling_burst_limit" {
  description = "Burst limit for throttling"
  type        = number
}

variable "throttling_rate_limit" {
  description = "Rate limit for throttling (requests per second)"
  type        = number
}

variable "quota_limit" {
  description = "Maximum number of requests allowed"
  type        = number
}

variable "quota_offset" {
  description = "Offset for the quota limit"
  type        = number
}

variable "quota_period" {
  description = "Time period for quota reset (e.g., DAY, WEEK, MONTH)"
  type        = string
}

variable "tags" {
  description = "Tags for resources"
  type        = map(string)
}
