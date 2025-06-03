//todo: we will add secondary indexes as we need them

resource "aws_dynamodb_table" "asset_table" {
  name         = "asset"
  billing_mode = "PAY_PER_REQUEST"

  attribute {
    name = "id"
    type = "N"
  }

  hash_key = "id"
}