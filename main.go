package main

import (
	"encoding/json"
	"github.com/astaxie/goorm"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
)

var orm ORM

func initDB() ORM {
	orm := goorm.NewORM(
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		"utf8")
	return orm
}

// Load .env variables
// Init db
func init() {
	gotenv.Load()
	orm := initDB()
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/charities", GetCharities).Methods("GET")
	router.HandleFunc("/charities/{id}", GetCharity).Methods("GET")
	port := getEnv("PORT", "8001")
	log.Println("Listening on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
