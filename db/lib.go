package db

import (
	"fmt"
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
func Connect() *gorm.DB {

	var (
		dbHost     = util.GetEnv("DB_HOST", "localhost")
		dbPort     = util.GetEnv("DB_PORT", "5432")
		dbUser     = util.GetEnv("DB_USER", "postgres")
		dbPassword = util.GetEnv("DB_PASSWORD", "")
		dbName     = util.GetEnv("DB_NAME", "development")
	)

	psqlConnectionString := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	for {
		var err error
		dbInstance, err = gorm.Open("postgres", psqlConnectionString)
		if err != nil {
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
