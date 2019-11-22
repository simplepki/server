data "archive_file" "user_get_token" {
  type = "zip"
  source_file = "${path.module}/../builds/user_get_token"
  output_path = "${path.module}/../builds/user_get_token.zip"
}
resource "aws_lambda_function" "user_get_token" {
	function_name = "user_get_token"
	role = "${aws_iam_role.user_get_token.arn}"
  filename = "${data.archive_file.user_get_token.output_path}"
  source_code_hash = "${data.archive_file.user_get_token.output_base64sha256}"
	handler = "user_get_token"
	runtime = "go1.x"
	timeout = "300"
	memory_size = 1024

  vpc_config {
    subnet_ids = split(",",var.lambda_subnets)

    security_group_ids = [
      "${aws_security_group.user_get_token.id}"
    ]
  }

  depends_on = ["aws_iam_role_policy_attachment.user_get_token_policy_attach_logs"]
}

resource "aws_security_group" "user_get_token" {
  name = "user_get_token"

   egress {
    from_port = 0
    to_port = 0
    protocol = -1
    cidr_blocks = [
      "0.0.0.0/0"
    ]
  }
}

resource "aws_iam_role" "user_get_token" {
  name = "user_get_token"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF

}

resource "aws_iam_role_policy_attachment" "user_get_token_policy_attach_logs" {
  role = "${aws_iam_role.user_get_token.name}"
  policy_arn = "${aws_iam_policy.logging_policy.arn}"
}

resource "aws_iam_role_policy_attachment" "user_get_token_policy_attach_jwt" {
  role = "${aws_iam_role.user_get_token.name}"
  policy_arn = "${aws_iam_policy.secrets_manager_jwt.arn}"
}

resource "aws_iam_role_policy_attachment" "user_get_token_policy_attach_ec2" {
  role = "${aws_iam_role.user_get_token.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}