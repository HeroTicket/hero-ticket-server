terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.25.0"
    }

    tls = {
      source  = "hashicorp/tls"
      version = "4.0.4"
    }
  }
}

provider "aws" {
  region = "ap-northeast-2"
}

provider "aws" {
  alias  = "us-east-1"
  region = "us-east-1"
}

provider "tls" {

}
