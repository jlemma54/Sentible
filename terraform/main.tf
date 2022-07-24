# Require TF version to be same as or greater than 0.12.13
terraform {
  required_version = ">=0.12.13"
  backend "s3" {
    bucket         = "reddit-sentiment-state-file"
    key            = "terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "reddit-sentiment-dynamo-db"
    encrypt        = true
  }
}

# Download any stable version in AWS provider of 2.36.0 'master test merge'
provider "aws" {
  region  = "us-east-1"
}

module "bootstrap" {
  source                      = "./modules/bootstrap"
  name_of_s3_bucket           = "reddit-sentiment-state-file"
  dynamo_db_table_name        = "reddit-sentiment-dynamo-db"
}

