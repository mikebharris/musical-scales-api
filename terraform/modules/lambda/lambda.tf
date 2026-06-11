resource "aws_iam_role" "api_iam_role" {
  name                  = "${var.product}-api-iam-role"
  force_detach_policies = true
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_cloudwatch_log_group" "api_cloudwatch_log_group" {
  name              = "/aws/lambda/${aws_lambda_function.api_lambda_function.function_name}"
  retention_in_days = 1
}

resource "aws_iam_role_policy_attachment" "api_policy_attachment_execution" {
  role       = aws_iam_role.api_iam_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

data "archive_file" "api_lambda_function_distribution" {
  source_file = "../lambdas/${var.product}/bootstrap"
  output_path = "../lambdas/${var.product}/${var.product}.zip"
  type        = "zip"
}

resource "aws_s3_object" "api_lambda_function_distribution_bucket_object" {
  bucket = var.distribution_bucket
  key    = "lambdas/${var.product}/${var.product}.zip"
  source = data.archive_file.api_lambda_function_distribution.output_path
  etag   = filemd5(data.archive_file.api_lambda_function_distribution.output_path)
}

resource "aws_lambda_function" "api_lambda_function" {
  function_name    = var.product
  role             = aws_iam_role.api_iam_role.arn
  handler          = "bootstrap"
  runtime          = "provided.al2023"
  architectures    = ["arm64"]
  s3_bucket        = aws_s3_object.api_lambda_function_distribution_bucket_object.bucket
  s3_key           = aws_s3_object.api_lambda_function_distribution_bucket_object.key
  source_code_hash = data.archive_file.api_lambda_function_distribution.output_md5
  timeout          = 15
  memory_size      = 128

  tags = {
    Name          = "${var.product}.lambda"
    Contact       = var.contact
    Project       = var.product
    Orchestration = var.orchestration
    Description   = "API for returning musical scales"
  }
}

resource "aws_lambda_function_url" "api_lambda_function_url" {
  authorization_type = "NONE"
  function_name      = aws_lambda_function.api_lambda_function.function_name

  cors {
    allow_credentials = false
    allow_origins     = ["*"]
    allow_methods     = ["GET"]
    allow_headers     = ["content-type"]
    expose_headers    = []
    max_age           = 86400
  }
}
