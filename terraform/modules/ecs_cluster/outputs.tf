output "cluster_name" {
  value = aws_ecs_cluster.hobby_streamer.name
}

output "cluster_arn" {
  value = aws_ecs_cluster.hobby_streamer.arn
}