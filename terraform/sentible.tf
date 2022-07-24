resource "aws_s3_bucket" "reddit-sentible-calc-bucket" {
  bucket = "reddit-sentible-calc-bucket"
  acl = "private"

  # Tells AWS to encrypt the S3 bucket at rest by default
  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }

  lifecycle {
    prevent_destroy = true
  }

  versioning {
    enabled = true
  }

  tags = {
    Terraform = "true"
  }
}

resource "aws_vpc_endpoint" "s3" {
    vpc_id = "${aws_vpc.sentible-prod-vpc.id}"
    service_name = "com.amazonaws.us-east-1.s3"
}

resource "aws_vpc_endpoint_route_table_association" "s3_vpc_association"{
    route_table_id = "${aws_route_table.sentible-prod-public-crt.id}"
    vpc_endpoint_id = "${aws_vpc_endpoint.s3.id}"
}

# resource "aws_kms_key" "rds-key" {
#     description = "key to encrypt rds password"
#   tags = {
#     Name = "my-rds-kms-key"
#   }
# }

# resource "aws_kms_alias" "rds-kms-alias" {
#   target_key_id = "${aws_kms_key.rds-key.id}"
#   name = "alias/rds-kms-key"
# }

data "aws_kms_secrets" "rds-secret" {
  secret {
    name = "sentibledb_password"
    payload = "AQICAHiYmLl2Zq7fkV/714vkkfzRIubNuUntQivfIrGzxM7CjgFZlWWkabU0HQf+cNNVS08QAAAAajBoBgkqhkiG9w0BBwagWzBZAgEAMFQGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMH3a0nDFBUSQk3+BUAgEQgCeipF61qaRAIlkdOQCg5xrJtmHPlhsNyxs1KAw1Hwgv1GG+DXgKECw="
  }
}

resource "aws_db_instance" "sentibledb" {
  allocated_storage           = 20
  storage_type                = "gp2"
  engine                      = "mysql"
  engine_version              = "5.7"
  instance_class              = "db.t3.micro"
  name                        = "sentibledb"
  username                    = "admin"
  password                    = "${data.aws_kms_secrets.rds-secret.plaintext["sentibledb_password"]}"
  parameter_group_name        = "default.mysql5.7"
  db_subnet_group_name        = "${aws_db_subnet_group.rds-private-subnet.name}"
  vpc_security_group_ids      = ["${aws_security_group.db-sg.id}"]
  allow_major_version_upgrade = true
  auto_minor_version_upgrade  = true
  backup_retention_period     = 35
  backup_window               = "22:00-23:00"
  maintenance_window          = "Sat:00:00-Sat:03:00"
  multi_az                    = true
  skip_final_snapshot         = true
  publicly_accessible         = true


  tags = {
      Name = "sent"
  }
}

resource "aws_kms_key" "sentible-rds-key" {
    description = "key to encrypt rds password"
  tags = {
    Name = "my-rds-kms-key"
  }
}

resource "aws_kms_alias" "sentible-rds-kms-alias" {
  target_key_id = "${aws_kms_key.sentible-rds-key.id}"
  name = "alias/sentible-rds-kms-key"
}

