variable "access_key" {}
variable "access_secret" {}
variable "key_name" {}
variable "private_key" {}
variable "db_password" {}

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

variable "db_user" {
  default = "postgres"
}

variable "db_port" {
  default = "5432"
}

variable "db_name" {
  default = "staging"
}

variable "cms_image" {
  default = "sihoang/charity-management-serv"
}

variable "cms_seeder_container" {
  default = "cms_seeder"
}

locals {
  common_tags = {
    Terraform   = "true"
    Environment = "${var.env}"
    App         = "${var.app_name}"
  }
}

provider "aws" {
  access_key = "${var.access_key}"
  secret_key = "${var.access_secret}"
  region     = "${var.region}"
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  name   = "${var.app_name}"

  cidr             = "172.31.0.0/16"
  public_subnets   = ["172.31.32.0/20", "172.31.64.0/20"]
  database_subnets = ["172.31.48.0/20", "172.31.80.0/20"]

  enable_nat_gateway   = true
  enable_dns_hostnames = true
  enable_dns_support   = true

  azs  = ["${var.region}a", "${var.region}b"]
  tags = "${local.common_tags}"
}

module "web_server_sg" {
  source = "terraform-aws-modules/security-group/aws"
  name   = "web-server-sg"

  # TODO make the port open within VPC behind the load balancer
  description = "Security group for web-server with HTTP ports open to public"
  vpc_id      = "${module.vpc.vpc_id}"

  egress_cidr_blocks  = ["0.0.0.0/0"]
  egress_rules        = ["all-tcp"]
  ingress_cidr_blocks = ["0.0.0.0/0"]
  ingress_rules       = ["https-443-tcp", "http-80-tcp", "ssh-tcp"]

  tags = "${local.common_tags}"
}

module "postgres_sg" {
  source      = "terraform-aws-modules/security-group/aws"
  name        = "postgres-sg"
  description = "Security group for postgres DB with port open within VPC"

  vpc_id = "${module.vpc.vpc_id}"

  computed_ingress_with_source_security_group_id = [
    {
      from_port                = "${var.db_port}"
      to_port                  = "${var.db_port}"
      protocol                 = "tcp"
      source_security_group_id = "${module.web_server_sg.this_security_group_id}"
    },
  ]

  number_of_computed_ingress_with_source_security_group_id = 1
}

resource "aws_instance" "master" {
  ami                         = "${var.ami_id}"
  associate_public_ip_address = true
  instance_type               = "t2.small"
  key_name                    = "${var.key_name}"
  monitoring                  = true
  vpc_security_group_ids      = ["${module.web_server_sg.this_security_group_id}"]
  subnet_id                   = "${module.vpc.public_subnets[0]}"
  tags                        = "${local.common_tags}"
}

resource "null_resource" "provision" {
  connection {
    host        = "${aws_instance.master.public_dns}"
    user        = "ubuntu"
    private_key = "${file(var.private_key)}"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo apt-get -y remove docker docker-engine docker.io",
      "sudo apt-get update",
      "sudo apt-get -y install apt-transport-https ca-certificates curl software-properties-common zip",
      "curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -",
      "sudo add-apt-repository \"deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable\"",
      "sudo apt-get update",
      "sudo apt-get -y install docker-ce",
      "sudo docker pull ${var.cms_image}",
      <<EOF
        sudo docker run \
          --restart always \
          --name cms \
          -p 80:8001 \
          -e DB_HOST=${module.db.this_db_instance_address} \
          -e DB_PORT=${var.db_port} \
          -e DB_PASSWORD=${var.db_password} \
          -e DB_NAME=${var.db_name} \
          -e DB_USER=${var.db_user} \
          -d ${var.cms_image}
      EOF
      ,
      "curl -O https://apps.irs.gov/pub/epostcard/data-download-pub78.zip",
      "unzip -o data-download-pub78.zip",
      "sudo docker stop ${var.cms_seeder_container} || true",
      "sudo docker rm ${var.cms_seeder_container} || true",
      "sudo docker run --name ${var.cms_seeder_container} -t -d ${var.cms_image} /bin/bash",
      <<EOF
        sudo docker cp data-download-pub78.txt \
          ${var.cms_seeder_container}:/go/src/github.com/WeTrustPlatform/charity-management-serv/seed/data.txt
      EOF
      ,
      <<EOF
        sudo docker exec \
          -e DB_HOST=${module.db.this_db_instance_address} \
          -e DB_PORT=${var.db_port} \
          -e DB_PASSWORD=${var.db_password} \
          -e DB_NAME=${var.db_name} \
          -e DB_USER=${var.db_user} \
          -d ${var.cms_seeder_container} make seeder
      EOF
      ,
    ]
  }
}

module "db" {
  source            = "terraform-aws-modules/rds/aws"
  identifier        = "${var.app_name}"
  engine            = "postgres"
  engine_version    = "10.4"
  instance_class    = "db.t2.small"
  allocated_storage = 10
  storage_type      = "gp2"

  maintenance_window = "Mon:00:00-Mon:02:00"
  backup_window      = "03:00-04:00"

  name     = "${var.db_name}"
  username = "${var.db_user}"
  password = "${var.db_password}"
  port     = "${var.db_port}"

  vpc_security_group_ids = ["${module.postgres_sg.this_security_group_id}"]
  subnet_ids             = "${module.vpc.database_subnets}"

  #deletion_protection = true

  family = "postgres10"
}

output "aws_instance_public_dns" {
  value = "${aws_instance.master.public_dns}"
}

output "db_host" {
  value = "${module.db.this_db_instance_address}"
}
