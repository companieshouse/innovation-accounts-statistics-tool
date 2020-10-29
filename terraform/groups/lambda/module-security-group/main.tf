resource "aws_security_group" "accounts_statistics_tool" {
  name        = "${var.environment}-${var.service}-lambda-into-vpc"
  description = "Outbound rules for accounts statistics tool lambda"
  vpc_id = var.vpc_id

  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_blocks     = ["0.0.0.0/0"]
  }
}

output "lambda_into_vpc_id" {
  value = aws_security_group.accounts_statistics_tool.id
}
