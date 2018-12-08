## Topology
# Public ==> CloudFront ==> ELB ==> EC2 ==> RDS
# CloudFront is not part of this config file
# ELB listen to 80 and 443. Public.
# EC2 only listen to 80 from ELB. Internal.
# RDS only listen to 5432 from EC2. Internal.

## Enforce SSL
# Option 1: Configure CloudFront
# Option 2: Configure nginx
# https://aws.amazon.com/premiumsupport/knowledge-center/redirect-http-https-elb

## Load balance
# TODO migrate to the next gen ALB
# TODO configure AWS ELB for multiple EC2's
# TODO configure nginx and docker for multiple api processes
# TODO configure RDS replication
# For now, we have:
# 1 instance behind ELB
# 1 api process
# CloudFront is manually configured to cache web static files

variable "access_key" {}
variable "access_secret" {}
variable "key_name" {}
variable "private_key" {}
variable "db_password" {}
variable "nginx_conf" {}
variable "ssl_certificate_id" {}

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
  default = "sihoang/charity-management-serv:latest"
}

variable "web_image" {
  default = "sihoang/charity-tcr:testnet-latest"
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

module "elb_sg" {
  source = "terraform-aws-modules/security-group/aws"
  name   = "elb-sg"

  description = "Security group for ELB with HTTP ports open to public"
  vpc_id      = "${module.vpc.vpc_id}"

  computed_egress_with_source_security_group_id = [
    {
      from_port                = "80"
      to_port                  = "80"
      protocol                 = "tcp"
      source_security_group_id = "${module.web_server_sg.this_security_group_id}"
    },
  ]

  number_of_computed_egress_with_source_security_group_id = 1

  ingress_cidr_blocks = ["0.0.0.0/0"]
  ingress_rules       = ["https-443-tcp", "http-80-tcp"]

  tags = "${local.common_tags}"
}

module "web_server_sg" {
  source = "terraform-aws-modules/security-group/aws"
  name   = "web-server-sg"

  # TODO make the port open within VPC behind the load balancer
  description = "Security group for web-server with HTTP ports open to public"
  vpc_id      = "${module.vpc.vpc_id}"

  egress_cidr_blocks = ["0.0.0.0/0"]
  egress_rules       = ["all-tcp"]

  ingress_cidr_blocks = ["0.0.0.0/0"]
  ingress_rules       = ["ssh-tcp"]

  computed_ingress_with_source_security_group_id = [
    {
      from_port                = "80"
      to_port                  = "80"
      protocol                 = "tcp"
      source_security_group_id = "${module.elb_sg.this_security_group_id}"
    },
  ]

  number_of_computed_ingress_with_source_security_group_id = 1

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

resource "null_resource" "provision_nginx" {
  connection {
    host        = "${aws_instance.master.public_dns}"
    user        = "ubuntu"
    private_key = "${file(var.private_key)}"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo apt-get -y install nginx",
      "sudo rm /etc/nginx/sites-enabled/* || true",
    ]
  }

  provisioner "file" {
    source      = "${var.nginx_conf}"
    destination = "./${var.nginx_conf}"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo mv ./${var.nginx_conf} /etc/nginx/sites-enabled/",
      "sudo service nginx reload",
    ]
  }
}

resource "null_resource" "provision_docker" {
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
      "curl -O https://apps.irs.gov/pub/epostcard/data-download-pub78.zip",
      "unzip -o data-download-pub78.zip",
    ]
  }
}

resource "null_resource" "provision_web" {
  depends_on = ["null_resource.provision_docker"]

  connection {
    host        = "${aws_instance.master.public_dns}"
    user        = "ubuntu"
    private_key = "${file(var.private_key)}"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo docker pull ${var.web_image}",
      "sudo docker stop web || true",
      "sudo docker rm web || true",
      <<EOF
        sudo docker run \
          --restart always \
          --name web \
          -p 8000:80 \
          -d ${var.web_image}
      EOF
      ,
    ]
  }
}

resource "null_resource" "provision_cms" {
  depends_on = ["null_resource.provision_docker"]

  connection {
    host        = "${aws_instance.master.public_dns}"
    user        = "ubuntu"
    private_key = "${file(var.private_key)}"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo docker pull ${var.cms_image}",
      "sudo docker stop cms || true",
      "sudo docker rm cms || true",
      <<EOF
        sudo docker run \
          --restart always \
          --name cms \
          -p 8001:8001 \
          -e DB_HOST=${module.db.this_db_instance_address} \
          -e DB_PORT=${var.db_port} \
          -e DB_PASSWORD=${var.db_password} \
          -e DB_NAME=${var.db_name} \
          -e DB_USER=${var.db_user} \
          -e ALLOWED_ORIGINS=* \
          -d ${var.cms_image}
      EOF
      ,
    ]
  }
}

resource "null_resource" "provision_seeder" {
  depends_on = ["null_resource.provision_docker", "null_resource.provision_cms"]

  connection {
    host        = "${aws_instance.master.public_dns}"
    user        = "ubuntu"
    private_key = "${file(var.private_key)}"
  }

  provisioner "remote-exec" {
    inline = [
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

resource "aws_elb" "elb" {
  name            = "${var.app_name}"
  security_groups = ["${module.elb_sg.this_security_group_id}"]
  subnets         = ["${module.vpc.public_subnets}"]

  listener {
    instance_port     = 80
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  listener {
    instance_port      = 80
    instance_protocol  = "http"
    lb_port            = 443
    lb_protocol        = "https"
    ssl_certificate_id = "${var.ssl_certificate_id}"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "HTTP:80/api/v0/charities"
    interval            = 30
  }

  instances                   = ["${aws_instance.master.id}"]
  cross_zone_load_balancing   = true
  idle_timeout                = 400
  connection_draining         = true
  connection_draining_timeout = 400
  tags                        = "${local.common_tags}"
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
  tags   = "${local.common_tags}"
}

output "aws_instance_public_dns" {
  value = "${aws_instance.master.public_dns}"
}

output "aws_elb_dns_name" {
  value = "${aws_elb.elb.dns_name}"
}

output "db_host" {
  value = "${module.db.this_db_instance_address}"
}
