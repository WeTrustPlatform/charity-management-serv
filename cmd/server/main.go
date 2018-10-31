package main

import (
	"log"
	"net/http"

	"github.com/WeTrustPlatform/charity-management-serv/db"
	"github.com/WeTrustPlatform/charity-management-serv/util"
	"github.com/gorilla/mux"
)

func main() {
	dbInstance := db.Connect()
	defer dbInstance.Close()
	dbInstance.AutoMigrate(&db.Charity{})

	router := mux.NewRouter()
	router.HandleFunc("/charities", db.GetCharities).Methods("GET")
	router.HandleFunc("/charities/{id}", db.GetCharity).Methods("GET")
	port := util.GetEnv("PORT", "8001")
	log.Println("Listening on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
