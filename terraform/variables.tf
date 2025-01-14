variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
}

variable "private_subnets" {
  description = "List of CIDR blocks for private subnets"
  type        = list(string)
}
