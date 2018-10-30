# Charity Management Serv
Work-in-progress


### Overview
Provide RESTful endpoints to access all the 501c3 organizations information.


### Getting started
- Install [go](https://golang.org/) and [dep](https://golang.github.io/dep/docs/installation.html) using methods of your choice.  They are available in most of Linux package repositories.
- Install [postgres](https://www.postgresql.org/download/).
- Install dependencies `make dep`.
- (Optional) Create your `.env` variables for your local configs. The default is:
```
export PORT=8001
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=
export DB_NAME=development
```
- (Optional) Seed database:
  * Download pub78 data at [irs](https://www.irs.gov/charities-non-profits/tax-exempt-organization-search-bulk-data-downloads).
  * Move the `.txt` to `seed/data.txt`.
  * `make seeder`.
- Launch the dev server `make server`.


### License
[GPL-3.0](https://www.gnu.org/licenses/gpl-3.0.txt) &copy; WeTrustPlatform
