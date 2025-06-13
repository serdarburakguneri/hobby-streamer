resource "aws_ecs_cluster" "hobby_streamer" {
  name = var.cluster_name

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = var.tags
}