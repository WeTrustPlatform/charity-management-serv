variable "access_key" {}
variable "access_secret" {}

variable "region" {
  default = "us-west-2"
}
# This ami_id depends on region
# Ubuntu 18.04 hvm:ebs-ssd
variable "ami_id" {
  default = "ami-079b4e9085609225c"
}
variable "env" {
  default = "staging"
}
variable "app_name" {
  default = "tcr-staging-1"
}
variable "subnet_id" {
  default = ""
}

locals {
  common_tags = {
    Terraform = "true"
    Environment = "${var.env}"
    App = "${var.app_name}"
  }
}

provider "aws" {
  access_key = "${var.access_key}"
  secret_key = "${var.access_secret}"
  region = "${var.region}"
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "${var.app_name}"
  cidr = "172.31.0.0/16"
  public_subnets = ["172.31.16.0/20"]

  enable_nat_gateway = true
  enable_dns_hostnames = true

  azs = ["${var.region}a"]
  tags = "${local.common_tags}"
}

module "web_server_sg" {
  source = "terraform-aws-modules/security-group/aws"

  name = "web-server-sg"
  description = "Security group for web-server with HTTP ports open within VPC"
  vpc_id = "${module.vpc.vpc_id}"

  egress_cidr_blocks = ["0.0.0.0/0"]
  egress_rules = ["all-tcp"]
  ingress_cidr_blocks = ["0.0.0.0/0"]
  ingress_rules = ["https-443-tcp", "http-80-tcp", "ssh-tcp"]

  tags = "${local.common_tags}"
}

module "ec2_instance" {
  source = "terraform-aws-modules/ec2-instance/aws"
  version = "1.12.0"

  name = "${var.app_name}"
  instance_count = 1

  ami = "${var.ami_id}"
  associate_public_ip_address = true
  instance_type = "t2.medium"
  key_name = "${var.app_name}"
  monitoring = true
  vpc_security_group_ids = ["${module.web_server_sg.this_security_group_id}"]
  subnet_id = "${module.vpc.public_subnets[0]}"

  tags = "${local.common_tags}"
}

output "public_dns" {
  value = "${module.ec2_instance.public_dns}"
}
