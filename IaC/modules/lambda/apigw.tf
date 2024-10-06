resource "aws_api_gateway_rest_api" "main" {
  name        = join("-", [var.account_name, local.region_short, "APIGW", var.project_name])
  description = "API Gateway for ${var.project_name}"
}

resource "aws_api_gateway_resource" "default" {
  rest_api_id = aws_api_gateway_rest_api.main.id
  parent_id   = aws_api_gateway_rest_api.main.root_resource_id
  path_part   = "$default"
}

resource "aws_api_gateway_method" "any" {
  rest_api_id   = aws_api_gateway_rest_api.main.id
  resource_id   = aws_api_gateway_resource.default.id
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "integration" {
  rest_api_id             = aws_api_gateway_rest_api.main.id
  resource_id             = aws_api_gateway_resource.default.id
  http_method             = aws_api_gateway_method.any.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.function.invoke_arn
}

resource "aws_api_gateway_domain_name" "api" {
  domain_name     = "api.${var.base_url}"
  certificate_arn = "arn:aws:acm:us-east-1:716182248480:certificate/4e6112a6-21da-4952-a44d-1df9d3a4c20f"
}

data "aws_route53_zone" "base_url" {
  name = var.base_url
}

resource "aws_route53_record" "api" {
  name    = aws_api_gateway_domain_name.api.domain_name
  type    = "A"
  zone_id = data.aws_route53_zone.base_url.zone_id

  alias {
    evaluate_target_health = true
    name                   = aws_api_gateway_domain_name.api.cloudfront_domain_name
    zone_id                = aws_api_gateway_domain_name.api.cloudfront_zone_id
  }
}

resource "aws_api_gateway_base_path_mapping" "api" {
  api_id      = aws_api_gateway_rest_api.main.id
  stage_name  = aws_api_gateway_stage.latest.stage_name
  domain_name = aws_api_gateway_domain_name.api.domain_name
}

# Lambda
resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.function.function_name
  principal     = "apigateway.amazonaws.com"

  # More: http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-control-access-using-iam-policies-to-invoke-api.html
  source_arn = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.main.id}/*/${aws_api_gateway_method.any.http_method}${aws_api_gateway_resource.default.path}"
}

# Deployment

resource "aws_api_gateway_deployment" "latest" {
  rest_api_id = aws_api_gateway_rest_api.main.id

  triggers = {
    redeployment = sha1(jsonencode(aws_api_gateway_rest_api.main.body))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "latest" {
  deployment_id = aws_api_gateway_deployment.latest.id
  rest_api_id   = aws_api_gateway_rest_api.main.id
  stage_name    = "latest"
}
