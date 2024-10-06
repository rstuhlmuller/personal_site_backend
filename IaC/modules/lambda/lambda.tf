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

resource "aws_iam_role" "iam_for_lambda" {
  name               = join("-", [var.account_name, local.region_short, "IAM-Role", var.project_name, "Lambda"])
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "null_resource" "go_build" {
  provisioner "local-exec" {
    command = "go build -tags lambda.norpc -o bootstrap ${var.path_to_go_file}"
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
}
