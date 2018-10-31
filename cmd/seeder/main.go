package main

import (
	"flag"
	"github.com/WeTrustPlatform/charity-management-serv/db"
	"github.com/WeTrustPlatform/charity-management-serv/seed"
)

func main() {
	var dryRun bool
	flag.BoolVar(&dryRun, "dryrun", false, "dryrun")
	flag.Parse()

	dbInstance := db.Connect()
	defer dbInstance.Close()
	dbInstance.AutoMigrate(&db.Charity{})

	seed.Populate(dbInstance, dryRun)
}
