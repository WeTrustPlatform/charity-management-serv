variable "remote_host" {}
variable "remote_user" {}
variable "private_key" {}
variable "nginx_conf" {}

variable "mainnet_image" {
  default = "wetrustplatform/staking-dapp:mainnet-latest"
}

variable "testnet_image" {
  default = "wetrustplatform/staking-dapp:testnet-latest"
}

variable "cms_image" {
  default = "sihoang/charity-management-serv:latest"
}

locals {
  connection = {
    host        = "${var.remote_host}"
    user        = "${var.remote_user}"
    private_key = "${file(var.private_key)}"
  }
}

resource "null_resource" "nginx" {
  // cannot re-use local.connection
  // https://github.com/hashicorp/terraform/issues/8616
  connection = {
    host        = "${local.connection["host"]}"
    user        = "${local.connection["user"]}"
    private_key = "${local.connection["private_key"]}"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo apt-get update",
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
      "sudo systemctl restart nginx",
    ]
  }
}

resource "null_resource" "docker" {
  depends_on = ["null_resource.nginx"]

  connection = {
    host        = "${local.connection["host"]}"
    user        = "${local.connection["user"]}"
    private_key = "${local.connection["private_key"]}"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo apt-get remove docker docker-engine docker.io containerd runc",
      "sudo apt-get update",
      "sudo apt-get install apt-transport-https ca-certificates curl gnupg2 software-properties-common",
      "curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add -",
      "sudo add-apt-repository \"deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable\"",
      "sudo apt-get update",
      "sudo apt-get install docker-ce docker-ce-cli containerd.io",
    ]
  }
}

resource "null_resource" "postgres" {
  depends_on = ["null_resource.docker"]

  connection = {
    host        = "${local.connection["host"]}"
    user        = "${local.connection["user"]}"
    private_key = "${local.connection["private_key"]}"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo docker pull postgres:10-alpine",
      <<EOF
        sudo docker run \
          --restart always \
          --name postgres \
          -e POSTGRES_DB=cms_development \
          -p 5432:5432 \
          -v ~/data:/var/lib/postgresql/data \
          -d postgres:10-alpine
      EOF
      ,
    ]
  }
}

resource "null_resource" "mainnet" {
  depends_on = ["null_resource.postgres"]

  connection = {
    host        = "${local.connection["host"]}"
    user        = "${local.connection["user"]}"
    private_key = "${local.connection["private_key"]}"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo docker pull ${var.mainnet_image}",
      "sudo docker stop mainnet || true",
      "sudo docker rm mainnet || true",
      <<EOF
        sudo docker run \
          --restart always \
          --name mainnet \
          -p 8000:80 \
          -d ${var.mainnet_image}
      EOF
      ,
    ]
  }
}

resource "null_resource" "testnet" {
  depends_on = ["null_resource.mainnet"]

  connection = {
    host        = "${local.connection["host"]}"
    user        = "${local.connection["user"]}"
    private_key = "${local.connection["private_key"]}"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo docker pull ${var.testnet_image}",
      "sudo docker stop testnet || true",
      "sudo docker rm testnet || true",
      <<EOF
        sudo docker run \
          --restart always \
          --name testnet \
          -p 7999:80 \
          -d ${var.testnet_image}
      EOF
      ,
    ]
  }
}

resource "null_resource" "cms" {
  depends_on = ["null_resource.testnet"]

  connection = {
    host        = "${local.connection["host"]}"
    user        = "${local.connection["user"]}"
    private_key = "${local.connection["private_key"]}"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo docker pull ${var.cms_image}",
      "sudo docker stop cms1 cms2 || true",
      "sudo docker rm cms1 cms2 || true",
      <<EOF
        sudo docker run \
          --restart always \
          --name cms1 \
          --network host \
          -e PORT=8001 \
          -p 8001:8001 \
          -v ~/:/workdir \
          -d ${var.cms_image}
      EOF
      ,
      <<EOF
        sudo docker run \
          --restart always \
          --name cms2 \
          --network host \
          -e PORT=8002 \
          -p 8002:8002 \
          -v ~/:/workdir \
          -d ${var.cms_image}
      EOF
      ,
    ]
  }
}
