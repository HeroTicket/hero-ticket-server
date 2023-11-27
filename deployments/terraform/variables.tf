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
      name_tag             = "Hero Ticket Server"
    }
    "hero-ticket-issuer-node" = {
      name                 = "hero-ticket-issuer-node"
      image_tag_mutability = "MUTABLE"
      scan_on_push         = true
      name_tag             = "Hero Ticket Issuer Node"
    }
    "hero-ticket-subscriber" = {
      name                 = "hero-ticket-subscriber"
      image_tag_mutability = "MUTABLE"
      scan_on_push         = true
      name_tag             = "Hero Ticket Subscriber"
    }
  }
}

variable "keypair_name" {
  type        = string
  description = "Name of the key pair to use for EC2 instances"
  default     = "hero-ticket-key-pair"
}

variable "vpc_cidr_block" {
  type        = string
  description = "CIDR block for the VPC"
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidr_blocks" {
  type        = list(string)
  description = "CIDR blocks for the public subnets"
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "private_subnet_cidr_blocks" {
  type        = list(string)
  description = "CIDR blocks for the private subnets"
  default     = ["10.0.3.0/24", "10.0.4.0/24"]
}

variable "availability_zones" {
  type        = list(string)
  description = "Availability zones to use for the subnets"
  default     = ["ap-northeast-2c", "ap-northeast-2a"]
}

variable "s3_bucket_name" {
  type        = string
  description = "Name of the S3 bucket to create"
  default     = "hero-ticket-bucket"
}

variable "manager_ip" {
  type        = string
  description = "IP address of the manager"
}
