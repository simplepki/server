resource "aws_iam_policy" "s3_policy" {
    name        = "s3_policy"
    description = "s3_policy"
    policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "0",
            "Effect": "Allow",
            "Action": [
                "s3:ListBucket",
                "s3:GetBucketLocation"
            ],
            "Resource": "arn:aws:s3:::${var.bucket}"
        },
        {
            "Sid": "1",
            "Effect": "Allow",
            "Action": "s3:*",
            "Resource": "arn:aws:s3:::${var.bucket}/*"
        }
    ]
}
EOF
}

# See also the following AWS managed policy: AWSLambdaBasicExecutionRole
resource "aws_iam_policy" "logging_policy" {
  name = "logging_policy"
  path = "/"
  description = "IAM policy for logging from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "certificate_ledger_access_policy" {
	name = "certificate_ledger_access_policy"
	description = "IAM policy for accessing dynamodb certificate ledger"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
	{
		"Effect": "Allow",
		"Action": [
			"dynamodb:BatchGetItem",
			"dynamodb:GetItem",
			"dynamodb:Query",
			"dynamodb:Scan",
			"dynamodb:BatchWriteItem",
			"dynamodb:PutItem",
			"dynamodb:UpdateItem"
		],
		"Resource": [
			"arn:aws:dynamodb:*:*:table/CertificateLedger",
			"arn:aws:dynamodb:*:*:table/CertificateLedger/*"
			]
	}
  ]
}
EOF
}

resource "aws_iam_policy" "secrets_manager_create_access_policy" {
	name = "secrets_manager_create_access_policy"
	description = "IAM policy for accessing secrets manager"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
	{
		"Effect": "Allow",
		"Action": [
			"secretsmanager:CreateSecret",
      "secretsmanager:DescribeSecret"
		],
		"Resource": "arn:aws:secretsmanager:*:*:*"
	}
  ]
}
EOF
}

resource "aws_iam_policy" "secrets_manager_get_access_policy" {
	name = "secrets_manager_get_access_policy"
	description = "IAM policy for accessing secrets manager"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
	{
		"Effect": "Allow",
		"Action": [
			"secretsmanager:GetSecretValue"
		],
		"Resource": "arn:aws:secretsmanager:*:*:*"
	}
  ]
}
EOF
}

resource "aws_iam_policy" "secrets_manager_mysql" {
	name = "secrets_manager_mysql"
	description = "IAM policy for accessing secrets manager"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
	{
		"Effect": "Allow",
		"Action": [
      "secretsmanager:ListSecret",
      "secretsmanager:DescribeSecret",
			"secretsmanager:GetSecretValue"
		],
		"Resource": "arn:aws:secretsmanager:us-west-1:953366327760:secret:mysql-WBluSP"
	}
  ]
}
EOF
}

resource "aws_iam_policy" "ec2_network_policy" {
  name = "ec2_network_policy"
  description= "IAM policy for accessing ec2 network interfaces"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents",
                "ec2:CreateNetworkInterface",
                "ec2:DescribeNetworkInterfaces",
                "ec2:DeleteNetworkInterface"
            ],
            "Resource": "*"
        }
    ]
}
EOF
}

data "aws_iam_policy" "VPCLambda" {
  arn = "arn:aws:iam::aws:policy/AWSLambdaVPCAccessExecutionRole"
} 