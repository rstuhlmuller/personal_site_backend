// any
resource "aws_api_gateway_resource" "proxy" {
  parent_id   = aws_api_gateway_rest_api.main.root_resource_id
  path_part   = "{proxy+}"
  rest_api_id = aws_api_gateway_rest_api.main.id

}
resource "aws_api_gateway_method" "any" {
  authorization = "NONE"
  http_method   = "ANY"
  resource_id   = aws_api_gateway_resource.proxy.id
  rest_api_id   = aws_api_gateway_rest_api.main.id
}

resource "aws_api_gateway_integration" "proxy" {
  http_method             = aws_api_gateway_method.any.http_method
  integration_http_method = "POST"
  resource_id             = aws_api_gateway_resource.proxy.id
  rest_api_id             = aws_api_gateway_rest_api.main.id
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.function.invoke_arn
}

resource "aws_api_gateway_method_response" "any" {
  http_method = aws_api_gateway_method.any.http_method
  resource_id = aws_api_gateway_method.any.resource_id
  rest_api_id = aws_api_gateway_rest_api.main.id
  status_code = "200"

  response_models = {
    "application/json" = "Empty"
  }

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = true
    "method.response.header.Access-Control-Allow-Methods" = true
    "method.response.header.Access-Control-Allow-Origin"  = true
  }
}

// visitor-count 
# resource "aws_api_gateway_resource" "visitor-count" {
#   parent_id   = aws_api_gateway_rest_api.main.root_resource_id
#   path_part   = "visitor_count"
#   rest_api_id = aws_api_gateway_rest_api.main.id
# }

# resource "aws_api_gateway_method" "visitor-count_GET" {
#   authorization = "NONE"
#   http_method   = "GET"
#   resource_id   = aws_api_gateway_resource.visitor-count.id
#   rest_api_id   = aws_api_gateway_rest_api.main.id
# }

# resource "aws_api_gateway_integration" "visitor-count_GET" {
#   http_method             = aws_api_gateway_method.visitor-count_GET.http_method
#   integration_http_method = "POST"
#   resource_id             = aws_api_gateway_resource.visitor-count.id
#   rest_api_id             = aws_api_gateway_rest_api.main.id
#   type                    = "AWS_PROXY"
#   uri                     = aws_lambda_function.function.invoke_arn
# }

# resource "aws_api_gateway_method_response" "visitor-count_GET" {
#   http_method = aws_api_gateway_method.visitor-count_GET.http_method
#   resource_id = aws_api_gateway_method.visitor-count_GET.resource_id
#   rest_api_id = aws_api_gateway_rest_api.main.id
#   status_code = "200"

#   response_models = {
#     "application/json" = "Empty"
#   }

#   response_parameters = {
#     "method.response.header.Access-Control-Allow-Headers" = true
#     "method.response.header.Access-Control-Allow-Methods" = true
#     "method.response.header.Access-Control-Allow-Origin"  = true
#   }
# }

# resource "aws_api_gateway_method_settings" "visitor-count_GET" {
#   method_path = join("/", [aws_api_gateway_resource.visitor-count.path, "GET"])
#   rest_api_id = aws_api_gateway_rest_api.main.id
#   stage_name  = aws_api_gateway_stage.latest.stage_name
#   settings {
#     throttling_burst_limit = 500
#     throttling_rate_limit  = 100
#   }
# }

# # resource "aws_api_gateway_method_response" "visitor-count_GET" {
# #   rest_api_id = aws_api_gateway_rest_api.main.id
# #   resource_id = aws_api_gateway_method.visitor-count_GET.resource_id
# #   http_method = aws_api_gateway_method.visitor-count_GET.http_method
# #   status_code = "200"
# # }

# resource "aws_api_gateway_method" "visitor-count_POST" {
#   authorization = "NONE"
#   http_method   = "POST"
#   resource_id   = aws_api_gateway_resource.visitor-count.id
#   rest_api_id   = aws_api_gateway_rest_api.main.id
# }

# resource "aws_api_gateway_integration" "visitor-count_POST" {
#   http_method             = aws_api_gateway_method.visitor-count_POST.http_method
#   integration_http_method = "POST"
#   resource_id             = aws_api_gateway_resource.visitor-count.id
#   rest_api_id             = aws_api_gateway_rest_api.main.id
#   type                    = "AWS_PROXY"
#   uri                     = aws_lambda_function.function.invoke_arn
# }

# resource "aws_api_gateway_method_response" "visitor-count_POST" {
#   http_method = aws_api_gateway_method.visitor-count_POST.http_method
#   resource_id = aws_api_gateway_method.visitor-count_POST.resource_id
#   rest_api_id = aws_api_gateway_rest_api.main.id
#   status_code = "200"

#   response_models = {
#     "application/json" = "Empty"
#   }

#   response_parameters = {
#     "method.response.header.Access-Control-Allow-Headers" = true
#     "method.response.header.Access-Control-Allow-Methods" = true
#     "method.response.header.Access-Control-Allow-Origin"  = true
#   }
# }

# resource "aws_api_gateway_method_settings" "visitor-count_POST" {
#   method_path = join("/", [aws_api_gateway_resource.visitor-count.path, "POST"])
#   rest_api_id = aws_api_gateway_rest_api.main.id
#   stage_name  = aws_api_gateway_stage.latest.stage_name
#   settings {
#     throttling_burst_limit = 500
#     throttling_rate_limit  = 100
#   }
# }

# # resource "aws_api_gateway_method_response" "visitor-count_POST" {
# #   rest_api_id = aws_api_gateway_rest_api.main.id
# #   resource_id = aws_api_gateway_method.visitor-count_POST.resource_id
# #   http_method = aws_api_gateway_method.visitor-count_POST.http_method
# #   status_code = "200"
# # }

# resource "aws_api_gateway_method" "visitor-count_OPTIONS" {
#   authorization = "NONE"
#   http_method   = "OPTIONS"
#   resource_id   = aws_api_gateway_resource.visitor-count.id
#   rest_api_id   = aws_api_gateway_rest_api.main.id
# }

# resource "aws_api_gateway_integration" "visitor-count_OPTIONS" {
#   http_method             = aws_api_gateway_method.visitor-count_OPTIONS.http_method
#   integration_http_method = "POST"
#   resource_id             = aws_api_gateway_resource.visitor-count.id
#   rest_api_id             = aws_api_gateway_rest_api.main.id
#   type                    = "AWS_PROXY"
#   uri                     = aws_lambda_function.function.invoke_arn
# }

# resource "aws_api_gateway_method_response" "visitor-count_OPTIONS" {
#   http_method = aws_api_gateway_method.visitor-count_OPTIONS.http_method
#   resource_id = aws_api_gateway_method.visitor-count_OPTIONS.resource_id
#   rest_api_id = aws_api_gateway_rest_api.main.id
#   status_code = "200"

#   response_models = {
#     "application/json" = "Empty"
#   }

#   response_parameters = {
#     "method.response.header.Access-Control-Allow-Headers" = true
#     "method.response.header.Access-Control-Allow-Methods" = true
#     "method.response.header.Access-Control-Allow-Origin"  = true
#   }
# }

