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
