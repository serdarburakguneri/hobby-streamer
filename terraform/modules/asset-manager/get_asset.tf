resource "aws_lambda_function" "get_asset" {
  function_name = "get_asset"
  role          = aws_iam_role.asset_lambda_role.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  filename      = "${path.module}/../build/get_asset.zip"
  source_code_hash = filebase64sha256("${path.module}/../build/get_asset.zip")

  environment {
    variables = {
      TABLE_NAME = aws_dynamodb_table.asset_table.name
    }
  }
}

resource "aws_api_gateway_resource" "asset_id" {
  rest_api_id = var.api_id
  parent_id   = aws_api_gateway_resource.assets.id
  path_part   = "{id}"
}

resource "aws_api_gateway_method" "get_asset" {
  rest_api_id      = var.api_id
  resource_id      = aws_api_gateway_resource.asset_id.id
  http_method      = "GET"
  authorization    = "NONE"
  api_key_required = true
}

resource "aws_api_gateway_integration" "get_asset" {
  rest_api_id             = var.api_id
  resource_id             = aws_api_gateway_resource.asset_id.id
  integration_http_method = "GET"
  http_method             = "GET"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.get_asset.invoke_arn
}

resource "aws_lambda_permission" "get_asset" {
  statement_id  = "AllowGetAssetInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.get_asset.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "arn:aws:execute-api:${var.aws_region}:${var.account_id}:${var.api_id}/${var.stage_name}/GET/assets/*"
}