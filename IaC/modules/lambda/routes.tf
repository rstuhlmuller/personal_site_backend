// root
resource "aws_api_gateway_method" "root_GET" {
  authorization = "NONE"
  http_method   = "GET"
  resource_id   = aws_api_gateway_rest_api.main.root_resource_id
  rest_api_id   = aws_api_gateway_rest_api.main.id
}

resource "aws_api_gateway_integration" "root_GET" {
  http_method             = aws_api_gateway_method.root_GET.http_method
  integration_http_method = "POST"
  resource_id             = aws_api_gateway_rest_api.main.root_resource_id
  rest_api_id             = aws_api_gateway_rest_api.main.id
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.function.invoke_arn
}

// ping 
resource "aws_api_gateway_resource" "ping" {
  parent_id   = aws_api_gateway_rest_api.main.root_resource_id
  path_part   = "ping"
  rest_api_id = aws_api_gateway_rest_api.main.id
}

resource "aws_api_gateway_method" "ping_GET" {
  authorization = "NONE"
  http_method   = "GET"
  resource_id   = aws_api_gateway_resource.ping.id
  rest_api_id   = aws_api_gateway_rest_api.main.id
}

resource "aws_api_gateway_integration" "ping_GET" {
  http_method             = aws_api_gateway_method.ping_GET.http_method
  integration_http_method = "POST"
  resource_id             = aws_api_gateway_resource.ping.id
  rest_api_id             = aws_api_gateway_rest_api.main.id
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.function.invoke_arn
}
