resource "aws_dynamodb_table" "basic-dynamodb-table" {
  name           = join("-", [var.account_name, local.region_short, "DDB", var.project_name])
  billing_mode   = "PAY_PER_REQUEST"
#   read_capacity  = 20
#   write_capacity = 20
  hash_key       = "ID"
  range_key      = "ProjectName"

  attribute {
    name = "ID"
    type = "S"
  }

  attribute {
    name = "ProjectName"
    type = "S"
  }

  #   attribute {
  #     name = "TopScore"
  #     type = "N"
  #   }

  #   ttl {
  #     attribute_name = "TimeToExist"
  #     enabled        = true
  #   }

  #   global_secondary_index {
  #     name               = "GameTitleIndex"
  #     hash_key           = "GameTitle"
  #     range_key          = "TopScore"
  #     write_capacity     = 10
  #     read_capacity      = 10
  #     projection_type    = "INCLUDE"
  #     non_key_attributes = ["UserId"]
  #   }

  tags = local.tags
}
