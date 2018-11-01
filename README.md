# Charity Management Serv
Work-in-progress


### Overview
Provide RESTful endpoints to access all the 501c3 organizations information.


### Getting started
- Install [go](https://golang.org/) and [dep](https://golang.github.io/dep/docs/installation.html) using methods of your choice.  They are available in most of Linux package repositories.
- Install [postgres](https://www.postgresql.org/download/).
- Install dependencies `make dep`.
- Build binary `make build`. All the binaries are in the auto-generated folder `bin/`.


### Making life easier
- Create your `.env` variables for your local configs. The default configs are:
```
# server port
export PORT=8001

# pagination
export PER_PAGE=10

# database connection
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=
export DB_NAME=development
```
- Seed database:
  * Download pub78 data at [irs](https://www.irs.gov/charities-non-profits/tax-exempt-organization-search-bulk-data-downloads).
  * Move the `.txt` to `seed/data.txt`.
  * Run `make seeder`.
- Live reload:
  * Install [fswatch](https://github.com/emcrisostomo/fswatch).
  * Start dev server `make dev`.
  * It will rebuild the server and restart when there are changes in `*.go`.


### Linting
- Install [gometalinter](https://github.com/alecthomas/gometalinter). Make sure you run `gometalinter --install` at least once.
- Run `make lint`
- If it fails and shows the list of files at the `Step: goimports`, then it is because the import sections of the files are not properly format.
- [vim-go](https://github.com/fatih/vim-go) runs `gofmt` on save by default. If you like it to run `goimports` instead, please refer to https://github.com/fatih/vim-go/issues/207


### License
[GPL-3.0](https://www.gnu.org/licenses/gpl-3.0.txt) &copy; WeTrustPlatform
