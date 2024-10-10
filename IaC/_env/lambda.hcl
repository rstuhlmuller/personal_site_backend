terraform {
  source = "${dirname(find_in_parent_folders())}/modules/lambda"
}

locals {
  common_vars  = yamldecode(file(find_in_parent_folders("common.yml")))
  account_vars = read_terragrunt_config(find_in_parent_folders("account.hcl")).locals
  region_vars  = read_terragrunt_config(find_in_parent_folders("region.hcl")).locals
}

dependency "dynamodb" {
  config_path = "${get_original_terragrunt_dir()}/../database"
}

inputs = {
  account_name    = local.account_vars.account_name
  region_short    = local.region_vars.region_short
  project_name    = local.common_vars.project_name
  path_to_go_file = "${dirname(find_in_parent_folders())}/../cmd/main.go"
  base_url        = local.account_vars.fqdn
  dynamodb_table  = dependency.dynamodb.outputs.dynamodb_table
  db_arn          = dependency.dynamodb.outputs.db_arn
}