resource "aws_lambda_function" "create_ca" {
	function_name = "create_ca"
	role = "${aws_iam_role.create_ca_role.arn}"
	s3_bucket = "${var.bucket}"
	s3_key = "create_ca_lambda.zip"
	handler = "builds/create_ca_lambda"
	runtime = "go1.x"
	timeout = "10"
	memory_size = 1024
	source_code_hash = "${filebase64sha256("../builds/create_ca_lambda.zip")}"

  vpc_config {
    subnet_ids = [
      "subnet-64506701"
    ]

    security_group_ids = [
      "${aws_security_group.create_ca_lambda.id}"
    ]
  }

  depends_on = ["aws_iam_role_policy_attachment.create_ca_role_policy_attach_logs"]
}

resource "aws_security_group" "create_ca_lambda" {
  name = "create_ca_lambda"

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

resource "aws_iam_role" "create_ca_role" {
  name = "create_ca_role"

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

resource "aws_iam_role_policy_attachment" "create_ca_role_policy_attach_s3" {
  role = "${aws_iam_role.create_ca_role.name}"
  policy_arn = "${aws_iam_policy.s3_policy.arn}"
}

resource "aws_iam_role_policy_attachment" "create_ca_role_policy_attach_logs" {
  role = "${aws_iam_role.create_ca_role.name}"
  policy_arn = "${aws_iam_policy.logging_policy.arn}"
}

resource "aws_iam_role_policy_attachment" "ca_role_policy_attach_secrets_manager_create" {
	role = "${aws_iam_role.create_ca_role.name}"
	policy_arn = "${aws_iam_policy.secrets_manager_create_access_policy.arn}"
}

resource "aws_iam_role_policy_attachment" "ca_role_policy_attach_mysql" {
  role = "${aws_iam_role.create_ca_role.name}"
  policy_arn = "${aws_iam_policy.secrets_manager_mysql.arn}"
}

resource "aws_iam_role_policy_attachment" "create_ca_role_policy_attach_ec2" {
  role = "${aws_iam_role.create_ca_role.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}