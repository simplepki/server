//not using api gateway but file name is fitting

resource "aws_lb" "gateway" {
    name = "gateway"
    load_balancer_type = "application"

    security_groups = [
        aws_security_group.allow_gateway.id
    ]

    subnets = split(",", var.gateway_subnets)


}

resource "aws_security_group" "allow_gateway" {
	name        = "allow_gateway"
	description = "allows http to gateway from current ip address"
  
	vpc_id = var.vpc_id

	ingress {
		from_port = 80
		to_port = 80
		protocol = "tcp"
		cidr_blocks = ["${chomp(data.http.myip.body)}/32"]
	}
	
	egress {
		from_port = 0
		to_port = 0
		protocol = "-1"
		cidr_blocks =[ "0.0.0.0/0"]
	}

}

resource "aws_lb_target_group" "lambda_router" {
  name        = "lambda-router"
  target_type = "lambda"
}

resource "aws_lambda_permission" "lambda_router" {
  statement_id  = "ALBRouterExecution"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.router.arn}"
  principal     = "elasticloadbalancing.amazonaws.com"
  source_arn    = "${aws_lb_target_group.lambda_router.arn}"
}

resource "aws_lb_target_group_attachment" "router" {
  target_group_arn = "${aws_lb_target_group.lambda_router.arn}"
  target_id        = "${aws_lambda_function.router.arn}"
  depends_on       = ["aws_lambda_permission.lambda_router"]
}

resource "aws_lb_listener" "router_listener" {
  load_balancer_arn = "${aws_lb.gateway.arn}"
  port              = "80"
  protocol          = "HTTP"
  //ssl_policy        = "ELBSecurityPolicy-2016-08"

  default_action {
    type             = "forward"
    target_group_arn = "${aws_lb_target_group.lambda_router.arn}"
  }
}

data "archive_file" "router" {
  type = "zip"
  source_file = "${path.module}/../builds/lambda_router"
  output_path = "${path.module}/../builds/lambda_router.zip"
}

resource "aws_lambda_function" "router" {
	function_name = "lambda-router"
	role = "${aws_iam_role.router_role.arn}"
	filename = "${data.archive_file.router.output_path}"
  source_code_hash = "${data.archive_file.router.output_base64sha256}"
	handler = "lambda_router"
	runtime = "go1.x"
	timeout = "300"
	memory_size = 1024

  depends_on = ["aws_iam_role_policy_attachment.router_policy_attach_logs"]
}

resource "aws_iam_role" "router_role" {
  name = "router_role"

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

resource "aws_iam_role_policy_attachment" "router_policy_attach_logs" {
  role = "${aws_iam_role.router_role.name}"
  policy_arn = "${aws_iam_policy.logging_policy.arn}"
}