[![Build Status](https://travis-ci.org/WeTrustPlatform/charity-management-serv.svg?branch=master)](https://travis-ci.org/WeTrustPlatform/charity-management-serv)

# Charity Management Serv
Work-in-progress


### Overview
Provide RESTful endpoints to access all the 501c3 organizations information and projects on spring.wetrust.io


### Getting started
- Install [go](https://golang.org/) and [dep](https://golang.github.io/dep/docs/installation.html) using methods of your choice.  They are available in most of Linux package repositories.
- Clone this repo to $GO_PATH/src/github.com/WeTrustPlatform/charity-management-serv
- Install dependencies `make dep`.
- Build binary `make build`. All the binaries are in the auto-generated folder `bin/`.
- Install and launch [postgres](https://www.postgresql.org/download/) or use Docker:
```
docker pull postgres:10-alpine
docker --rm -p 5432:5432 -e POSTGRES_DB=cms_development postgres:10-alpine
```
- Launch server
```
make launch
```


### Making dev life easier
- Create your `.env` variables for your local configs. The default configs are:
```
# server port
PORT=8001

# pagination
PER_PAGE=10

# database connection
DATABASE_URL="postgres://postgres:@localhost:5432/cms_development?sslmode=disable"
```
- Seed database:
  * Prefer the existing `seed/data_test.txt`. Real data can be downloaded at [irs](https://www.irs.gov/charities-non-profits/tax-exempt-organization-search-bulk-data-downloads).
  * Run `data=seed/data_test.txt make seeder`.
- Live reload: It will auto build and reload the server as you change source code.
  * Install [fswatch](https://github.com/emcrisostomo/fswatch).
  * Start dev server `make dev`.
  * It will rebuild and restart the server when there are changes in `*.go`.


### Docker
If you would like to use this repo as a dependency and do not care to modify the code, then you can get it up running quickly by using [docker-compose](https://docs.docker.com/compose/).
  * Launch `docker-compose up`
  * Seed `docker run -it --rm --network host -v $(pwd)/seed:/seed charity-management-serv_api cms-seeder -data /seed/data_test.txt`


### Linting
- Install [gometalinter](https://github.com/alecthomas/gometalinter). Make sure you run `gometalinter --install` at least once.
- Run `make lint`
- If it fails and shows the list of files at the `Step: goimports`, then it is because the import sections of the files are not properly format.
- [vim-go](https://github.com/fatih/vim-go) runs `gofmt` on save by default. If you like it to run `goimports` instead, please refer to https://github.com/fatih/vim-go/issues/207


### Setting up infrastructure
- Use [terraform](https://www.terraform.io/) to provise the AWS infra and deploy docker containers.
- Go to the deployment env i.e. `cd deployment/staging`
- Configure your `terraform.tfvars` as below:
```
access_key =  <IAM user with permissions to create EC2, VPC, RDS and CloudFront(optional). Otherwise, you'll get 403 Forbidden>

access_secret = <secret key when creating new IAM user>

key_name = <the pem file's name>

private_key = </path/to/pem - for example: /home/user/.ssh/key_name.pem>

db_password = <postgres password>

nginx_conf = "./nginx.conf"

ssl_certificate_id = <get cert for your domain via AWS ACM. Make sure it in the same region. For example: arn:aws:acm:us-west-2:xxxxxx:certificate/xxxxxxxxx"
```
- Other default configs are declared at the beginning of `servers.tf`
- Provise infra and deploy for the first time:
```
terraform apply
```


### Release
This app uses Docker and Terraform for deployment. The images need to be created and published first.

- To create and publish latest image:
```
./docker.sh
```

- To create, publish image and tag release:
```
./docker.sh v1.1.0
```

- Manually push git tags to remote:
```
git push origin,upstream --tags
```

- To deploy the latest docker images to staging:
```
cd deployment/staging
terraform taint null_resource.provision_cms // this is for charity-management-serv
terraform taint null_resource.provision_web // this is for staking-dapp front-end

terraform apply
```


### License
[GPL-3.0](https://www.gnu.org/licenses/gpl-3.0.txt) &copy; WeTrustPlatform
