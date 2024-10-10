data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role_policy" "dynamodb" {
  name = join("-", [var.account_name, local.region_short, "IAM-Role-Policy", var.project_name, "DynamoDB"])

  role   = aws_iam_role.iam_for_lambda.id
  policy = data.aws_iam_policy_document.dynamodb.json
}

data "aws_iam_policy_document" "dynamodb" {
  statement {
    effect = "Allow"

    actions = [
      "dynamodb:GetItem",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem",
      "dynamodb:DeleteItem"
    ]
    resources = [
      var.db_arn
    ]
  }
}

resource "aws_iam_role" "iam_for_lambda" {
  name                = join("-", [var.account_name, local.region_short, "IAM-Role", var.project_name, "Lambda"])
  assume_role_policy  = data.aws_iam_policy_document.assume_role.json
  managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"]
}

resource "null_resource" "go_build" {
  triggers = {
    "code_sha" = sha1(join("", [for f in fileset(path.cwd, "*") : filesha1("${path.cwd}/${f}")])) # TODO: fix with actual go files
  }
  provisioner "local-exec" {
    command = "GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap ${var.path_to_go_file}"
  }
}

data "archive_file" "lambda" {
  type        = "zip"
  source_file = "${path.cwd}/bootstrap"
  output_path = "lambda_function_payload.zip"
  depends_on  = [null_resource.go_build]
}

resource "aws_lambda_function" "function" {
  # If the file is not in the current working directory you will need to include a
  # path.module in the filename.
  filename      = "lambda_function_payload.zip"
  function_name = join("-", [var.account_name, local.region_short, "Lambda", var.project_name])
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "bootstrap"

  source_code_hash = data.archive_file.lambda.output_base64sha256

  runtime = "provided.al2023"

  environment {
    variables = {
      DYNAMODB_TABLE = var.dynamodb_table
    }
  }
}
