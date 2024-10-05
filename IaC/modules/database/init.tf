terraform {
  required_version = "~> 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

locals {
  region_short = var.region_short
  tags = merge(
    var.tags,
    {
      Module = "rstuhlmuller/personal_site_backend/modules/database"
    }
  )
}
