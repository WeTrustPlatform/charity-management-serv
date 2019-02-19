[![Build Status](https://travis-ci.org/WeTrustPlatform/charity-management-serv.svg?branch=master)](https://travis-ci.org/WeTrustPlatform/charity-management-serv)

# Charity Management Serv
Work-in-progress


### Overview
Provide RESTful endpoints to access all the 501c3 organizations information and projects on spring.wetrust.io


### Getting started
- Install [go](https://golang.org/) and [dep](https://golang.github.io/dep/docs/installation.html) using methods of your choice.  They are available in most of Linux package repositories.
- Install [postgres](https://www.postgresql.org/download/) or use Docker:
```
docker pull postgres:10-alpine
docker --rm -p 5432:5432 -e POSTGRES_DB=cms_development postgres:10-alpine
```
- Install dependencies `make dep`.
- Build binary `make build`. All the binaries are in the auto-generated folder `bin/`.


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
  * Download pub78 data at [irs](https://www.irs.gov/charities-non-profits/tax-exempt-organization-search-bulk-data-downloads).
  * Move the `.txt` to `seed/data.txt`.
  * Run `data=data.txt make seeder`.
- Live reload: It will auto build and reload the server as you change source code.
  * Install [fswatch](https://github.com/emcrisostomo/fswatch).
  * Start dev server `make dev`.
  * It will rebuild and restart the server when there are changes in `*.go`.


### Docker
If you would like to use this repo as a dependency and do not care to modify the code, then you can get it up running quickly by using [docker-compose](https://docs.docker.com/compose/).
  * Make sure you have the DB variables in the `.env` as above.
  * Make sure you have `seed/data.txt` as above.
  * Launch `docker-compose up`
  * Seed `docker exec -e data=data.txt charity-management-serv_api_1 make seeder`


### Linting
- Install [gometalinter](https://github.com/alecthomas/gometalinter). Make sure you run `gometalinter --install` at least once.
- Run `make lint`
- If it fails and shows the list of files at the `Step: goimports`, then it is because the import sections of the files are not properly format.
- [vim-go](https://github.com/fatih/vim-go) runs `gofmt` on save by default. If you like it to run `goimports` instead, please refer to https://github.com/fatih/vim-go/issues/207


### Deployment
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
- To deploy the latest containers:
```
terraform taint null_resource.provision_cms // this is for charity-management-serv
terraform taint null_resource.provision_web // this is for staking-dapp front-end

terraform apply
```

### License
[GPL-3.0](https://www.gnu.org/licenses/gpl-3.0.txt) &copy; WeTrustPlatform
