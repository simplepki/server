data "http" "myip" {
  url = "http://ipv4.icanhazip.com"
}

variable "account" {
	type = string
}

variable "region" {
	type = string
}

variable "vpc_id" {
	type = string
}

variable "lambda_subnets" {
	type = string
}

variable "gateway_subnets" {
	type = string
}