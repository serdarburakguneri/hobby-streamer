provider "aws" {
  region  = "eu-north-1"
  profile = "sandbox"
}

data "aws_caller_identity" "current" {}

module "vpc" {
  source          = "./modules/vpc"
  vpc_cidr        = var.vpc_cidr
  private_subnets = var.private_subnets
}

module "storage-api" {
  source = "./modules/apigw"

  api_name               = var.storage_api_name
  api_description        = var.storage_api_description
  api_key_name           = var.storage_api_key_name
  api_key_description    = var.storage_api_key_description
  usage_plan_name        = var.storage_api_usage_plan_name
  usage_plan_description = var.storage_api_usage_plan_description
  throttling_burst_limit = var.storage_api_throttling_burst_limit
  throttling_rate_limit  = var.storage_api_throttling_rate_limit
  quota_limit            = var.storage_api_quota_limit
  quota_offset           = var.storage_api_quota_offset
  quota_period           = var.storage_api_quota_period
  stage_name             = var.storage_api_stage_name
  tags                   = var.tags

  depends_on = [module.vpc]
}

module "storage" {
  source               = "./modules/storage"
  account_id           = data.aws_caller_identity.current.account_id
  api_id               = module.storage-api.api_id
  api_root_resource_id = module.storage-api.root_resource_id
  tags                 = var.tags
  aws_region           = var.aws_region
  s3_bucket_name       = var.s3_bucket_name
  stage_name           = var.storage_api_stage_name

  depends_on = [module.vpc, module.storage-api]
}



