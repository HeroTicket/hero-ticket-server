output "ecr_repo_names" {
  value = {
    for k, v in module.ecr : k => v.ecr_repo_name
  }
}

output "ecr_repo_arns" {
  value = {
    for k, v in module.ecr : k => v.ecr_repo_arn
  }
}

output "ecr_repo_registry_ids" {
  value = {
    for k, v in module.ecr : k => v.ecr_repo_registry_id
  }
}

output "ecr_repo_repository_url" {
  value = {
    for k, v in module.ecr : k => v.ecr_repo_repository_url
  }
}

output "hero_ticket_issuer_node_public_ips" {
  value = data.aws_instances.hero_ticket_issuer_node_instances.public_ips
}

output "hero_ticket_server_public_ips" {
  value = data.aws_instances.hero_ticket_server_instances.public_ips
}
