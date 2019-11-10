data "archive_file" "cert_create_intermediate" {
  type = "zip"
  source_file = "${path.module}/../builds/cert_create_intermediate"
  output_path = "${path.module}/../builds/cert_create_intermediate.zip"
}
resource "aws_lambda_function" "create_intermediate" {
	function_name = "cert_create_intermediate"
	role = "${aws_iam_role.create_intermediate_role.arn}"
	filename = "${data.archive_file.cert_create_intermediate.output_path}"
  source_code_hash = "${data.archive_file.cert_create_intermediate.output_base64sha256}"
	handler = "cert_create_intermediate"
	runtime = "go1.x"
	timeout = "300"
	memory_size = 1024

  vpc_config {
    subnet_ids = [
      "subnet-64506701"
    ]

    security_group_ids = [
      "${aws_security_group.intermediate_lambda.id}"
    ]
  }

  depends_on = ["aws_iam_role_policy_attachment.create_intermediate_role_policy_attach_logs"]
}

resource "aws_security_group" "intermediate_lambda" {
  name = "create_intermediate_lambda"

  /*egress {
    from_port = 3306
    to_port = 3306
    protocol = "tcp"
    cidr_blocks = [
      "172.31.0.0/16"
    ]
  }

  egress {
    from_port = 443
    to_port = 443
    protocol = "tcp"
    cidr_blocks = [
      "0.0.0.0/0"
    ]
  }

  egress {
    from_port = 80
    to_port = 80
    protocol = "tcp"
    cidr_blocks = [
      "0.0.0.0/0"
    ]
  }*/

   egress {
    from_port = 0
    to_port = 0
    protocol = -1
    cidr_blocks = [
      "0.0.0.0/0"
    ]
  }
}

resource "aws_iam_role" "create_intermediate_role" {
  name = "create_intermediate_role"

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

/*resource "aws_iam_role_policy_attachment" "create_intermediate_role_policy_attach_s3" {
  role = "${aws_iam_role.create_intermediate_role.name}"
  policy_arn = "${aws_iam_policy.s3_policy.arn}"
}*/

resource "aws_iam_role_policy_attachment" "create_intermediate_role_policy_attach_logs" {
  role = "${aws_iam_role.create_intermediate_role.name}"
  policy_arn = "${aws_iam_policy.logging_policy.arn}"
}

resource "aws_iam_role_policy_attachment" "intermediate_role_policy_attach_secrets_manager_get" {
	role = "${aws_iam_role.create_intermediate_role.name}"
	policy_arn = "${aws_iam_policy.secrets_manager_get_cert_access_policy.arn}"
}

resource "aws_iam_role_policy_attachment" "intermediate_role_policy_attach_secrets_manager_create" {
	role = "${aws_iam_role.create_intermediate_role.name}"
	policy_arn = "${aws_iam_policy.secrets_manager_create_cert_access_policy.arn}"
}

resource "aws_iam_role_policy_attachment" "intermediate_role_policy_attach_mysql" {
  role = "${aws_iam_role.create_intermediate_role.name}"
  policy_arn = "${aws_iam_policy.secrets_manager_mysql.arn}"
}

resource "aws_iam_role_policy_attachment" "intermediate_role_policy_attach_ec2" {
  role = "${aws_iam_role.create_intermediate_role.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}