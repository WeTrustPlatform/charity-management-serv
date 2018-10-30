package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
)

// Load .env variables
func init() {
	gotenv.Load()
}

func connectDB() *gorm.DB {
	var (
		dbHost     = getEnv("DB_HOST", "localhost")
		dbPort     = getEnv("DB_PORT", "5432")
		dbUser     = getEnv("DB_USER", "postgres")
		dbPassword = getEnv("DB_PASSWORD", "")
		dbName     = getEnv("DB_NAME", "development")
	)

	psqlConnectionString := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := gorm.Open("postgres", psqlConnectionString)
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

var db *gorm.DB
var err error

func main() {
	db = connectDB()
	defer db.Close()
	db.AutoMigrate(&Charity{})

	router := mux.NewRouter()
	router.HandleFunc("/charities", GetCharities).Methods("GET")
	router.HandleFunc("/charities/{id}", GetCharity).Methods("GET")
	port := getEnv("PORT", "8001")
	log.Println("Listening on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
