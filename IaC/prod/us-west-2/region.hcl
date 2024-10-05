locals {
  region       = "us-west-2"
  region_short = "USW2"

  default_tags = {
    Region = local.region
  }
}