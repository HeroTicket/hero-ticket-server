variable "common_tags" {
  type        = map(string)
  description = "Common tags to apply to all resources"
  default = {
    "Project"   = "Hero Ticket"
    "ManagedBy" = "Terraform"
  }
}

variable "ecr_repos" {
  type = map(object({
    name                 = string
    image_tag_mutability = string
    scan_on_push         = bool
    name_tag             = string
  }))
  description = "ECR repositories to create"
  default = {
    "hero-ticket-server" = {
      name                 = "hero-ticket-server"
      image_tag_mutability = "MUTABLE"
      scan_on_push         = true
      name_tag             = "hero-ticket-server"
    }
    "hero-ticket-swagger" = {
      name                 = "hero-ticket-swagger"
      image_tag_mutability = "MUTABLE"
      scan_on_push         = true
      name_tag             = "hero-ticket-swagger"
    }
  }
}
