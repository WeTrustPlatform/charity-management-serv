package seed

import (
	"github.com/WeTrustPlatform/charity-management-serv/db"
	"github.com/jinzhu/gorm"
)

// Populate will upload all the records in pub78 to DB
func Populate(dbInstance *gorm.DB) {
	// Stubbed for now
	// TODO load data from pub78 text
	charity1 := db.Charity{
		Name:              "Open Library",
		City:              "San Francisco",
		State:             "CA",
		Country:           "United States",
		EIN:               "94-3242767",
		DeductibilityCode: "PC",
		Website:           "https://openlibrary.org",
	}

	charity2 := db.Charity{
		Name:              "Code for America",
		City:              "San Francisco",
		State:             "CA",
		Country:           "United States",
		EIN:               "27-1067272",
		DeductibilityCode: "PC",
		Website:           "https://www.codeforamerica.org",
	}

	dbInstance.Create(&charity1)
	dbInstance.Create(&charity2)
}
