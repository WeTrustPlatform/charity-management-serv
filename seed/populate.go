package seed

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/WeTrustPlatform/charity-management-serv/db"
	"github.com/jinzhu/gorm"
)

// PopulateCausesFromSpring will call Spring causes' endpoint
// If the record already exist,
// it will update the IsOnSpring to true
func PopulateCausesFromSpring(dbInstance *gorm.DB, dryRun bool) {
	resp, requestError := http.Get("http://spring.wetrust.io/api/v1/causes")
	if requestError != nil {
		panic(requestError)
	}
	defer resp.Body.Close()

	body, readingError := ioutil.ReadAll(resp.Body)
	if readingError != nil {
		panic(readingError)
	}

	result := []SpringCause{}
	jsonError := json.Unmarshal(body, &result)
	if jsonError != nil {
		panic(jsonError)
	}

	for index, value := range result {
		if len(value.Name) == 0 || len(value.StakingID) == 0 {
			fmt.Printf("Invalid record at %d. Move on...", index)
			continue
		}

		charity := db.Charity{
			StakingID:  value.StakingID,
			Name:       value.Name,
			IsOnSpring: true,
		}

		if !dryRun {
			var (
				whereCondition = db.Charity{StakingID: charity.StakingID}
				updatedValue   = charity
			)

			dbInstance.Where(whereCondition).Assign(updatedValue).FirstOrCreate(&charity)
		} else {
			fmt.Println("Dryrun mode. Do nothing.")
		}

		fmt.Println(charity)
	}
}

// Populate501c3FromIRS will insert all the records in pub78 to DB.
// If the record already exists,
// it will update DB with new values in the pub78 txt
func Populate501c3FromIRS(dbInstance *gorm.DB, file string, dryRun bool) {
	fileWithAbsolutePath, err := filepath.Abs(file)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(filepath.Clean(fileWithAbsolutePath))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lines, err := NewCSVReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	for index, line := range lines {
		charity, err := ParseCharity(line)
		if err != nil || charity == nil {
			fmt.Printf("Invalid record at %d. Move on...", index)
			continue
		}

		if !dryRun {
			var (
				whereCondition = db.Charity{StakingID: charity.StakingID}
				// all the required fields available in pub78
				updatedValue = db.Charity{
					Name:              charity.Name,
					City:              charity.City,
					State:             charity.State,
					Country:           charity.Country,
					StakingID:         charity.StakingID,
					DeductibilityCode: charity.DeductibilityCode,
					Is501c3:           true,
				}
			)

			dbInstance.Where(whereCondition).Assign(updatedValue).FirstOrCreate(&charity)
		} else {
			fmt.Println("Dryrun mode. Do nothing.")
		}

		// just print out the charity struct for dryRun
		// TODO provide more info on how DB would change
		// like new records, updated fields, ...
		fmt.Println(charity)
	}
}

// Populate gets data from Spring and IRS
func Populate(dbInstance *gorm.DB, filename string, dryRun bool) {
	PopulateCausesFromSpring(dbInstance, dryRun)
	Populate501c3FromIRS(dbInstance, filename, dryRun)
}
