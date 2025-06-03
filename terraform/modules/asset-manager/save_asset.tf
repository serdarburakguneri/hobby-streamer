resource "aws_lambda_function" "save_asset" {
  function_name = "save_asset"
  role          = aws_iam_role.asset_lambda_role.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  filename      = "${path.module}/../build/save_asset.zip"
  source_code_hash = filebase64sha256("${path.module}/../build/save_asset.zip")

  environment {
    variables = {
      TABLE_NAME = aws_dynamodb_table.asset_table.name
    }
  }
}

resource "aws_api_gateway_resource" "assets" {
  rest_api_id = var.api_id
  parent_id   = var.api_root_resource_id
  path_part   = "assets"
}

resource "aws_api_gateway_method" "post_asset" {
  rest_api_id      = var.api_id
  resource_id      = aws_api_gateway_resource.assets.id
  http_method      = "POST"
  authorization    = "NONE"
  api_key_required = true
}

resource "aws_api_gateway_integration" "save_asset" {
  rest_api_id             = var.api_id
  resource_id             = aws_api_gateway_resource.assets.id
  integration_http_method = "POST"
  http_method             = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.save_asset.invoke_arn
}

resource "aws_lambda_permission" "save_asset" {
  statement_id  = "AllowSaveAssetInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.save_asset.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "arn:aws:execute-api:${var.aws_region}:${var.account_id}:${var.api_id}/${var.stage_name}/POST/assets"
}