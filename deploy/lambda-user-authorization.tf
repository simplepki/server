data "archive_file" "user_authorization" {
  type = "zip"
  source_file = "${path.module}/../builds/user_authorization"
  output_path = "${path.module}/../builds/user_authorization.zip"
}
resource "aws_lambda_function" "user_authorization" {
	function_name = "user_authorization"
	role = "${aws_iam_role.user_authorization.arn}"
  filename = "${data.archive_file.user_authorization.output_path}"
  source_code_hash = "${data.archive_file.user_authorization.output_base64sha256}"
	handler = "user_authorization"
	runtime = "go1.x"
	timeout = "300"
	memory_size = 1024

  vpc_config {
    subnet_ids = [
      "subnet-64506701"
    ]

    security_group_ids = [
      "${aws_security_group.user_authorization.id}"
    ]
  }

  depends_on = ["aws_iam_role_policy_attachment.user_authorization_policy_attach_logs"]
}

resource "aws_security_group" "user_authorization" {
  name = "user_authorization"

   egress {
    from_port = 0
    to_port = 0
    protocol = -1
    cidr_blocks = [
      "0.0.0.0/0"
    ]
  }
}

resource "aws_iam_role" "user_authorization" {
  name = "user_authorization"

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

resource "aws_iam_role_policy_attachment" "user_authorization_policy_attach_logs" {
  role = "${aws_iam_role.user_authorization.name}"
  policy_arn = "${aws_iam_policy.logging_policy.arn}"
}

resource "aws_iam_role_policy_attachment" "user_authorization_policy_attach_jwt" {
  role = "${aws_iam_role.user_authorization.name}"
  policy_arn = "${aws_iam_policy.secrets_manager_jwt.arn}"
}

resource "aws_iam_role_policy_attachment" "user_authorization_policy_attach_ec2" {
  role = "${aws_iam_role.user_authorization.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}