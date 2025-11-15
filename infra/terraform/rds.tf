resource "aws_db_subnet_group" "main" {
  name       = "resume-app-subnets"
  subnet_ids = var.private_subnet_ids
}

resource "aws_db_instance" "postgres" {
  identifier              = "resume-app-db"
  engine                  = "postgres"
  engine_version          = "16.3"
  instance_class          = "db.t3.micro"
  allocated_storage       = 20
  username                = var.db_username
  password                = var.db_password
  db_name                 = var.db_name
  skip_final_snapshot     = true
  vpc_security_group_ids  = [var.db_security_group_id]
  db_subnet_group_name    = aws_db_subnet_group.main.name
  publicly_accessible     = false
}

variable "private_subnet_ids" { type = list(string) }
variable "db_security_group_id" { type = string }
variable "db_username" { type = string default = "postgres" }
variable "db_password" { type = string default = "changeMe123!" }
variable "db_name" { type = string default = "resumeapp" }

output "rds_endpoint" { value = aws_db_instance.postgres.address }
