variable "account_name" {
  type        = string
  description = "Name of the account"
}

variable "region_short" {
  type        = string
  description = "Short name of the region"
}

variable "project_name" {
  type        = string
  description = "Name of the project"
}

variable "tags" {
  type        = map(any)
  description = "Tags"
  default     = {}
}

variable "path_to_go_file" {
  type        = string
  description = "Path to go file"
}

variable "base_url" {
  type        = string
  description = "Base URL"
}

variable "dynamodb_table" {
  type        = string
  description = "Name of the DynamoDB table"
}
