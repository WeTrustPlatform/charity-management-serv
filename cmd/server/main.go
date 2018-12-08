package main

import (
	"log"
	"net/http"

	"github.com/WeTrustPlatform/charity-management-serv/db"
	"github.com/WeTrustPlatform/charity-management-serv/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
)

// Router set up mux Router handlers
func Router() *mux.Router {
	root := "/api/v0"
	router := mux.NewRouter()
	router.HandleFunc(root+"/charities", db.GetCharities).Methods("GET")
	router.HandleFunc(root+"/charities/{id}", db.GetCharity).Methods("GET")
	return router
}

// DB initialize and connect to DB based on .env
func DB() *gorm.DB {
	dbInstance := db.Connect()
	return dbInstance
}

// CorsPolicy specify policy
func CorsPolicy() *cors.Cors {
	corsPolicy := cors.New(cors.Options{
		AllowedOrigins: []string{util.GetEnv("ALLOWED_ORIGINS", "http://localhost:8000")},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodOptions,
		},
		AllowedHeaders: []string{"*"},
	})
	return corsPolicy
}

func main() {
	dbInstance := DB()
	defer dbInstance.Close()

	router := Router()
	port := util.GetEnv("PORT", "8001")
	log.Println("Listening on http://localhost:" + port)

	corsPolicy := CorsPolicy()
	handler := corsPolicy.Handler(router)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
