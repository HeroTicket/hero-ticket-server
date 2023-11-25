module "ecr" {
  source = "./modules/ecr"

  for_each = var.ecr_repos

  name                 = each.value.name
  image_tag_mutability = each.value.image_tag_mutability
  scan_on_push         = each.value.scan_on_push
  common_tags          = merge(var.common_tags, { "Name" = each.value.name_tag })
}
