package seed

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/WeTrustPlatform/charity-management-serv/db"
	"github.com/jinzhu/gorm"
	"os"
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
func NewCSVReader(file *os.File) *csv.Reader {
	reader := csv.NewReader(file)
	reader.Comma = '|'
	return reader
}

// Populate will upload all the records in pub78 to DB
func Populate(dbInstance *gorm.DB, persist bool) {

	// Load data from txt file
	filename := "seed/data.txt"
	f, err := os.Open(filename)
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

		fmt.Println(charity)

		if persist {
			dbInstance.Create(charity)
		}
	}
}
