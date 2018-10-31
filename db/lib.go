package db

import (
	"fmt"

	"github.com/WeTrustPlatform/charity-management-serv/util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // wrap postgres driver
)

var dbInstance *gorm.DB
var err error

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

	dbInstance, err = gorm.Open("postgres", psqlConnectionString)
	if err != nil {
		panic("failed to connect database")
	}
	return dbInstance
}
