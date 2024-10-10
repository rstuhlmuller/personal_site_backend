output "dynamodb_table" {
  value = aws_dynamodb_table.basic-dynamodb-table.name
}

output "db_arn" {
  value = aws_dynamodb_table.basic-dynamodb-table.arn
}
