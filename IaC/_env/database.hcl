terraform {
  source = "${dirname(find_in_parent_folders())}/modules/database"
}

locals {
  common_vars  = yamldecode(file(find_in_parent_folders("common.yml")))
  account_vars = read_terragrunt_config(find_in_parent_folders("account.hcl")).locals
  region_vars  = read_terragrunt_config(find_in_parent_folders("region.hcl")).locals
}

inputs = {
  account_name = local.account_vars.account_name
  region_short = local.region_vars.region_short
  project_name = local.common_vars.project_name
}