data "http" "myip" {
  url = "http://ipv4.icanhazip.com"
}

variable "bastion_enabled" {
	type = number
	default = 1
}

variable "bastion_subnet" {
	type = string
}

resource "aws_security_group" "allow_bastion_ssh" {
	name        = "allow_bastion_ssh"
	description = "allows ssh to bastion from current ip address"
  
	vpc_id = var.vpc_id

	ingress {
		from_port = 22
		to_port = 22
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

data "aws_ami" "amazon-linux-2-ami" {
	most_recent = true
	
	owners = [
		"137112412989"
	]
	
	filter {
		name   = "owner-alias"
		values = ["amazon"]
	}
	
	filter {
		name   = "name"
		values = ["amzn2-ami-hvm*"]
	}
	
	filter {
		name   = "architecture"
		values = ["x86_64"]
	}
}


resource "aws_instance" "bastion" {
	count = var.bastion_enabled
	
	ami           = "${data.aws_ami.amazon-linux-2-ami.id}"
	instance_type = "t2.micro"
	
	associate_public_ip_address = true
	
	subnet_id = var.bastion_subnet

	vpc_security_group_ids = [
		aws_security_group.allow_bastion_ssh.id
	]
}

output "bastion_instance_id" {
	value = "${aws_instance.bastion.*.id}"
}