variable "name" {
  description = "Name of the ECS service (container, task, log prefix)"
  type        = string
}

variable "image" {
  description = "Docker image to deploy"
  type        = string
}

variable "cpu" {
  description = "Task CPU units"
  type        = number
  default     = 256
}

variable "memory" {
  description = "Task memory in MiB"
  type        = number
  default     = 512
}

variable "cluster_arn" {
  description = "ARN of the ECS cluster"
  type        = string
}

variable "container_port" {
  description = "Port on which the container listens"
  type        = number
}

variable "desired_count" {
  description = "Number of desired running tasks"
  type        = number
  default     = 1
}

variable "execution_role_arn" {
  type        = string
  description = "IAM role for ECS to pull container/logs"
}

variable "task_role_arn" {
  type        = string
  description = "IAM role for the task to assume"
}

variable "subnet_ids" {
  type        = list(string)
  description = "List of subnet IDs for Fargate tasks"
}

variable "security_group_id" {
  type        = string
  description = "Security group for Fargate tasks"
}

variable "assign_public_ip" {
  type        = bool
  description = "Whether to assign public IP"
  default     = false
}

variable "listener_arn" {
  type        = string
  description = "ALB listener ARN"
}

variable "listener_priority" {
  type        = number
  description = "Priority for listener rule"
  default     = 100
}

variable "target_group_arn" {
  type        = string
  description = "ALB target group ARN"
}

variable "path_patterns" {
  type        = list(string)
  description = "List of URL path patterns to match in ALB rule"
}

variable "region" {
  type        = string
  description = "AWS region (for logs)"
}

variable "environment" {
  type        = map(string)
  description = "Environment variables for container"
  default     = {}
}