provider "aws" {
  region  = "eu-north-1"
  profile = "sandbox"
}

module "vpc" {
  source          = "./modules/vpc"
  vpc_cidr        = var.vpc_cidr
  private_subnets = var.private_subnets
}
