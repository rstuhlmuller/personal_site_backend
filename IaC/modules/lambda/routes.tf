// root
# resource "aws_api_gateway_resource" "root" {
#   parent_id   = aws_api_gateway_rest_api.main.root_resource_id
#   path_part   = "{proxy+}"
#   rest_api_id = aws_api_gateway_rest_api.main.id
# }

# resource "aws_api_gateway_method" "root" {
#   authorization = "NONE"
#   http_method   = "ANY"
#   resource_id   = aws_api_gateway_resource.root.id
#   rest_api_id   = aws_api_gateway_rest_api.main.id
# }

# resource "aws_api_gateway_integration" "root" {
#   http_method             = aws_api_gateway_method.root.http_method
#   integration_http_method = "POST"
#   resource_id             = aws_api_gateway_resource.root.id
#   rest_api_id             = aws_api_gateway_rest_api.main.id
#   type                    = "AWS_PROXY"
#   uri                     = aws_lambda_function.function.invoke_arn
# }

// visitor-count 
resource "aws_api_gateway_resource" "visitor-count" {
  parent_id   = aws_api_gateway_rest_api.main.root_resource_id
  path_part   = "visitor-count"
  rest_api_id = aws_api_gateway_rest_api.main.id
}

resource "aws_api_gateway_method" "visitor-count_GET" {
  authorization = "NONE"
  http_method   = "GET"
  resource_id   = aws_api_gateway_resource.visitor-count.id
  rest_api_id   = aws_api_gateway_rest_api.main.id
}

resource "aws_api_gateway_integration" "visitor-count_GET" {
  http_method             = aws_api_gateway_method.visitor-count_GET.http_method
  integration_http_method = "POST"
  resource_id             = aws_api_gateway_resource.visitor-count.id
  rest_api_id             = aws_api_gateway_rest_api.main.id
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.function.invoke_arn
}

resource "aws_api_gateway_method" "visitor-count_POST" {
  authorization = "NONE"
  http_method   = "POST"
  resource_id   = aws_api_gateway_resource.visitor-count.id
  rest_api_id   = aws_api_gateway_rest_api.main.id
}

resource "aws_api_gateway_integration" "visitor-count_POST" {
  http_method             = aws_api_gateway_method.visitor-count_POST.http_method
  integration_http_method = "POST"
  resource_id             = aws_api_gateway_resource.visitor-count.id
  rest_api_id             = aws_api_gateway_rest_api.main.id
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.function.invoke_arn
}
