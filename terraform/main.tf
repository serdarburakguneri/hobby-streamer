provider "aws" {
  region  = "eu-north-1"
  profile = "sandbox"
}

data "aws_caller_identity" "current" {}

module "vpc" {
  source = "./modules/vpc"

  vpc_cidr        = var.vpc_cidr
  private_subnets = var.private_subnets
}

module "hobby-streamer-api" {
  source = "./modules/apigw"

  api_name               = var.api_name
  api_description        = var.api_description
  api_key_name           = var.api_key_name
  api_key_description    = var.api_key_description
  usage_plan_name        = var.api_usage_plan_name
  usage_plan_description = var.api_usage_plan_description
  throttling_burst_limit = var.api_throttling_burst_limit
  throttling_rate_limit  = var.api_throttling_rate_limit
  quota_limit            = var.api_quota_limit
  quota_offset           = var.api_quota_offset
  quota_period           = var.api_quota_period
  stage_name             = var.api_stage_name
  tags                   = var.tags

  depends_on = [module.vpc]
}

module "transcoder" {
  source = "./modules/transcoder"

  account_id                   = data.aws_caller_identity.current.account_id
  aws_region                   = var.aws_region
  api_id                       = module.hobby-streamer-api.api_id
  api_root_resource_id         = module.hobby-streamer-api.root_resource_id
  stage_name                   = var.api_stage_name
  sqs_queue_name               = var.sqs_queue_name
  sqs_queue_visibility_timeout = var.sqs_queue_visibility_timeout
  raw_storage_bucket_arn       = module.storage.raw_storage_bucket_arn
  tags                         = var.tags

  depends_on = [module.hobby-streamer-api]
}

module "storage" {
  source = "./modules/storage"

  account_id                        = data.aws_caller_identity.current.account_id
  api_id                            = module.hobby-streamer-api.api_id
  api_root_resource_id              = module.hobby-streamer-api.root_resource_id
  transcoding_queue_arn             = module.transcoder.transcoding_queue_arn
  tags                              = var.tags
  aws_region                        = var.aws_region
  raw_storage_s3_bucket_name        = var.raw_storage_s3_bucket_name
  stage_name                        = var.api_stage_name

  depends_on = [module.hobby-streamer-api]
}

module "triggers" {
  source = "./modules/triggers"

  raw_storage_bucket_id = module.storage.raw_storage_bucket_id
  transcoding_queue_arn = module.transcoder.transcoding_queue_arn

  depends_on = [module.storage, module.transcoder]
}



