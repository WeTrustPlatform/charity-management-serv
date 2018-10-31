package main

import (
	"github.com/WeTrustPlatform/charity-management-serv/db"
	"github.com/WeTrustPlatform/charity-management-serv/seed"
)

func main() {
	dbInstance := db.Connect()
	defer dbInstance.Close()
	dbInstance.AutoMigrate(&db.Charity{})

	seed.Populate(dbInstance, false)
}
