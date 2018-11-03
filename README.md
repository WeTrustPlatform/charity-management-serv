# Charity Management Serv
Work-in-progress


### Overview
Provide RESTful endpoints to access all the 501c3 organizations information.


### Getting started
- Install [go](https://golang.org/) and [dep](https://golang.github.io/dep/docs/installation.html) using methods of your choice.  They are available in most of Linux package repositories.
- Install [postgres](https://www.postgresql.org/download/).
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
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME=development
```
- Seed database:
  * Download pub78 data at [irs](https://www.irs.gov/charities-non-profits/tax-exempt-organization-search-bulk-data-downloads).
  * Move the `.txt` to `seed/data.txt`.
  * Run `make seeder`.
- Live reload: It will auto build and reload the server as you change source code.
  * Install [fswatch](https://github.com/emcrisostomo/fswatch).
  * Start dev server `make dev`.
  * It will rebuild the server and restart when there are changes in `*.go`.


### Docker
If you would like to use this repo as a dependency and do not care to modify the code, then you can get it up running quickly by using [docker-compose](https://docs.docker.com/compose/).
  * Make sure you have the DB variables in the `.env` as above.
  * Make sure you have `seed/data.txt` as above.
  * Launch `docker-compose up`
  * Seed `docker exec -it charity-management-serv_api_1 make seeder`


### Linting
- Install [gometalinter](https://github.com/alecthomas/gometalinter). Make sure you run `gometalinter --install` at least once.
- Run `make lint`
- If it fails and shows the list of files at the `Step: goimports`, then it is because the import sections of the files are not properly format.
- [vim-go](https://github.com/fatih/vim-go) runs `gofmt` on save by default. If you like it to run `goimports` instead, please refer to https://github.com/fatih/vim-go/issues/207


### License
[GPL-3.0](https://www.gnu.org/licenses/gpl-3.0.txt) &copy; WeTrustPlatform
