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

module "ecs_cluster" {
  source = "./modules/ecs_cluster"

  name = "hobby-streamer-cluster"
  tags = var.tags
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

module "storage" {
  source = "./modules/storage"

  account_id                        = data.aws_caller_identity.current.account_id
  api_id                            = module.hobby-streamer-api.api_id
  api_root_resource_id              = module.hobby-streamer-api.root_resource_id
  tags                              = var.tags
  aws_region                        = var.aws_region
  raw_storage_s3_bucket_name        = var.raw_storage_s3_bucket_name
  transcoded_storage_s3_bucket_name = var.transcoded_storage_s3_bucket_name
  thumbnail_storage_s3_bucket_name  = var.thumbnail_storage_s3_bucket_name
  stage_name                        = var.api_stage_name

  depends_on = [module.hobby-streamer-api]
}

module "asset_manager_service" {
  source = "./modules/fargate_service"

  name               = "asset-manager"
  image              = var.asset_manager_image
  region             = var.aws_region
  cluster_arn        = module.ecs_cluster.cluster_arn
  container_port     = 8080

  subnet_ids         = module.vpc.private_subnets
  security_group_id  = module.vpc.fargate_sg_id
  assign_public_ip   = false

  execution_role_arn = module.iam.execution_role_arn
  task_role_arn      = module.iam.task_role_arn

  listener_arn       = module.alb.listener_arn
  listener_priority  = 20
  target_group_arn   = module.alb.asset_manager_target_group_arn
  path_patterns      = ["/assets*", "/health"]

  environment = {
    PORT = "8080"
  }

  depends_on = [module.ecs_cluster]
}

module "transcoder_service" {
  source = "./modules/fargate_service"

  name               = "transcoder"
  image              = var.transcoder_image
  region             = var.aws_region
  cluster_arn        = module.ecs_cluster.cluster_arn
  container_port     = 8080

  subnet_ids         = module.vpc.private_subnets
  security_group_id  = module.vpc.fargate_sg_id
  assign_public_ip   = false

  execution_role_arn = module.iam.execution_role_arn
  task_role_arn      = module.iam.task_role_arn

  listener_arn       = module.alb.listener_arn
  listener_priority  = 30
  target_group_arn   = module.alb.transcoder_target_group_arn
  path_patterns      = ["/transcoder*", "/health"]

  environment = {
    PORT = "8080"
  }

  depends_on = [module.ecs_cluster]
}

module "events" {
  source = "./modules/events"

  raw_storage_bucket_id = module.storage.raw_storage_bucket_id
  transcoding_queue_arn = module.transcoder_service.environment["SQS_QUEUE_ARN"]

  depends_on = [module.storage]
}