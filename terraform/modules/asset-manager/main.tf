resource "aws_dynamodb_table" "metadata_table" {
  name         = "metadata"
  billing_mode = "PAY_PER_REQUEST"

  attribute {
    name = "id"
    type = "N"
  }

  hash_key = "id"
}

