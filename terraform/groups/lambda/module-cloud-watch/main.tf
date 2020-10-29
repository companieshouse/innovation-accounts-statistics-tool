resource "aws_cloudwatch_event_rule" "accounts_statistics_tool" {
  name        = "${var.service}-${var.environment}"
  description = "Call accounts statistics tool lambda"
  schedule_expression = var.cron_schedule
}

resource "aws_cloudwatch_event_target" "call_accounts_statistics_tool_lambda" {
    rule = aws_cloudwatch_event_rule.accounts_statistics_tool.name
    target_id = "${var.service}-${var.environment}"
    arn = var.lambda_arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_accounts_statistics_tool" {
    statement_id = "AllowExecutionFromCloudWatch"
    action = "lambda:InvokeFunction"
    function_name = "${var.service}-${var.environment}"
    principal = "events.amazonaws.com"
    source_arn = aws_cloudwatch_event_rule.accounts_statistics_tool.arn
}
