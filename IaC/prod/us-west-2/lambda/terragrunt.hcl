include "root" {
  path = find_in_parent_folders()
}

include "database" {
  path = "${dirname(find_in_parent_folders())}/_env/lambda.hcl"
}