package db

import (
	"log"
	"time"

	"github.com/WeTrustPlatform/charity-management-serv/util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // wrap postgres driver
)

var dbInstance *gorm.DB

// Connect initializes a DB connection based on .env configs
// default is postgres://postgres:@localhost:5432/development
// and return the DB instance
func Connect(retry bool) *gorm.DB {

	psqlConnectionString := util.GetEnv("DATABASE_URL", "postgres://postgres:@localhost:5432/cms_development?sslmode=disable")

	for {
		var err error
		dbInstance, err = gorm.Open("postgres", psqlConnectionString)
		if err != nil {
			if !retry {
				panic(err)
			}
			log.Println(err)
			time.Sleep(time.Duration(3) * time.Second)
			continue
		}

		break
	}

	dbInstance.AutoMigrate(&Charity{})

	// Custom queries those are not supported by Gorm
	dbInstance.Exec("ALTER TABLE charities ADD COLUMN IF NOT EXISTS tsv tsvector;")
	dbInstance.Exec(`
		UPDATE charities
		SET tsv = setweight(to_tsvector(name), 'A')
		|| setweight(to_tsvector(ein), 'B')
		|| setweight(to_tsvector(city), 'C')
		|| setweight(to_tsvector(state), 'D')
		;
		`)
	dbInstance.Exec("CREATE INDEX IF NOT EXISTS ix_charities_tsv ON charities USING GIN(tsv);")
	return dbInstance
}
