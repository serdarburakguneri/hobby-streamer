resource "aws_lambda_function" "list_assets" {
  function_name = "list_assets"
  role          = aws_iam_role.asset_lambda_role.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  filename      = "${path.module}/../build/list_assets.zip"
  source_code_hash = filebase64sha256("${path.module}/../build/list_assets.zip")

  environment {
    variables = {
      TABLE_NAME = aws_dynamodb_table.asset_table.name
    }
  }
}

resource "aws_api_gateway_method" "list_assets" {
  rest_api_id      = var.api_id
  resource_id      = aws_api_gateway_resource.assets.id
  http_method      = "GET"
  authorization    = "NONE"
  api_key_required = true
}

resource "aws_api_gateway_integration" "list_assets" {
  rest_api_id             = var.api_id
  resource_id             = aws_api_gateway_resource.assets.id
  integration_http_method = "GET"
  http_method             = "GET"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.list_assets.invoke_arn
}

resource "aws_lambda_permission" "list_assets" {
  statement_id  = "AllowListAssetsInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.list_assets.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "arn:aws:execute-api:${var.aws_region}:${var.account_id}:${var.api_id}/${var.stage_name}/GET/assets"
}