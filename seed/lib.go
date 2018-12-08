package seed

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/WeTrustPlatform/charity-management-serv/db"
	"github.com/jinzhu/gorm"
)

// ParseCharity constructs Charity from an array of string
// array has to conform to data order in the pub78 txt
func ParseCharity(line []string) (*db.Charity, error) {
	if len(line) < 6 {
		return nil, errors.New("Invalid record")
	}

	charity := db.Charity{
		EIN:               line[0],
		Name:              line[1],
		City:              line[2],
		State:             line[3],
		Country:           line[4],
		DeductibilityCode: line[5],
	}

	return &charity, nil
}

// NewCSVReader takes in a filename and returns a csv.Reader
func NewCSVReader(file io.Reader) *csv.Reader {
	reader := csv.NewReader(file)
	reader.Comma = '|'
	return reader
}

// Populate will insert all the records in pub78 to DB.
// If the record already exists,
// it will update DB with new values in the pub78 txt
func Populate(dbInstance *gorm.DB, filename string, dryRun bool) {
	// data file must be in the same folder as seed/lib.go
	_, thisFile, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(thisFile)
	filePath := filepath.Join(basePath, filename)
	f, err := os.Open(filepath.Clean(filePath))
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
				whereCondition = db.Charity{EIN: charity.EIN}
				// all the required fields available in pub78
				updatedValue = db.Charity{
					Name:              charity.Name,
					City:              charity.City,
					State:             charity.State,
					Country:           charity.Country,
					EIN:               charity.EIN,
					DeductibilityCode: charity.DeductibilityCode,
				}
			)

			dbInstance.Where(whereCondition).Assign(updatedValue).FirstOrCreate(&charity)
		}

		// just print out the charity struct for dryRun
		// TODO provide more info on how DB would change
		// like new records, updated fields, ...
		fmt.Println(charity)
	}
}
