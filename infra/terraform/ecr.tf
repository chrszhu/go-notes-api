resource "aws_ecr_repository" "app" {
  name                 = var.app_name
  image_tag_mutability = "MUTABLE"
  lifecycle_policy {
    policy = <<EOF
{
  "rules": [
    {
      "rulePriority": 1,
      "description": "Expire untagged images after 30 days",
      "selection": {
        "tagStatus": "untagged",
        "countType": "sinceImagePushed",
        "countUnit": "days",
        "countNumber": 30
      },
      "action": {"type": "expire"}
    }
  ]
}
EOF
  }
}

output "ecr_repository_url" {
  value = aws_ecr_repository.app.repository_url
}

variable "app_name" {
  type        = string
  description = "App / ECR repository name"
  default     = "resume-app"
}
