provider "aws" {
  region  = "eu-north-1"
  profile = "sandbox"
}

module "vpc" {
  source          = "./modules/vpc"
  vpc_cidr        = var.vpc_cidr
  private_subnets = var.private_subnets
}

module "apigw" {
  source = "./modules/apigw"

  api_name               = var.api_name
  api_key_name           = var.api_key_name
  api_key_description    = var.api_key_description
  usage_plan_name        = var.usage_plan_name
  usage_plan_description = var.usage_plan_description
  throttling_burst_limit = var.throttling_burst_limit
  throttling_rate_limit  = var.throttling_rate_limit
  quota_limit            = var.quota_limit
  quota_offset           = var.quota_offset
  quota_period           = var.quota_period
  tags                   = var.tags
}

module "storage" {
  source = "./modules/storage"

  s3_bucket_name = var.s3_bucket_name
  aws_region     = var.aws_region
  tags           = var.tags
}


