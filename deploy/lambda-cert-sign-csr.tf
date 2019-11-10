data "archive_file" "cert_sign_csr" {
  type = "zip"
  source_file = "${path.module}/../builds/cert_sign_csr"
  output_path = "${path.module}/../builds/cert_sign_csr.zip"
}
resource "aws_lambda_function" "sign_user_certificate" {
	function_name = "cert_sign_csr"
	role = "${aws_iam_role.sign_user_certificate_role.arn}"
  filename = "${data.archive_file.cert_sign_csr.output_path}"
  source_code_hash = "${data.archive_file.cert_sign_csr.output_base64sha256}"
	handler = "cert_sign_csr"
	runtime = "go1.x"
	timeout = "300"
	memory_size = 1024

  vpc_config {
    subnet_ids = [
      "subnet-64506701"
    ]

    security_group_ids = [
      "${aws_security_group.cert_sign_csr.id}"
    ]
  }

  depends_on = ["aws_iam_role_policy_attachment.sign_user_certificate_role_policy_attach_logs"]
}

resource "aws_security_group" "cert_sign_csr" {
  name = "cert_sign_csr"

   egress {
    from_port = 0
    to_port = 0
    protocol = -1
    cidr_blocks = [
      "0.0.0.0/0"
    ]
  }
}

resource "aws_iam_role" "sign_user_certificate_role" {
  name = "sign_user_certificate_role"

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

resource "aws_iam_role_policy_attachment" "sign_user_certificate_role_policy_attach_s3" {
  role = "${aws_iam_role.sign_user_certificate_role.name}"
  policy_arn = "${aws_iam_policy.s3_policy.arn}"
}

resource "aws_iam_role_policy_attachment" "sign_user_certificate_role_policy_attach_logs" {
  role = "${aws_iam_role.sign_user_certificate_role.name}"
  policy_arn = "${aws_iam_policy.logging_policy.arn}"
}

/*resource "aws_iam_role_policy_attachment" "sign_user_certificate_role_policy_attach_dynamodb" {
	role = "${aws_iam_role.sign_user_certificate_role.name}"
	policy_arn = "${aws_iam_policy.certificate_ledger_access_policy.arn}"
}*/

resource "aws_iam_role_policy_attachment" "sign_user_certificate_role_policy_attach_secrets_manager_get" {
	role = "${aws_iam_role.sign_user_certificate_role.name}"
	policy_arn = "${aws_iam_policy.secrets_manager_get_cert_access_policy.arn}"
}


resource "aws_iam_role_policy_attachment" "sign_user_certificate_role_policy_attach_mysql" {
  role = "${aws_iam_role.sign_user_certificate_role.name}"
  policy_arn = "${aws_iam_policy.secrets_manager_mysql.arn}"
}

resource "aws_iam_role_policy_attachment" "sign_user_certificate_role_policy_attach_ec2" {
  role = "${aws_iam_role.sign_user_certificate_role.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}