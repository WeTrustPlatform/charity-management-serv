package seed

import (
	"encoding/csv"
	"errors"
	"io"

	"github.com/WeTrustPlatform/charity-management-serv/db"
)

// SpringCause Cause object from spring.wetrust.io
type SpringCause struct {
	Name            string `json:"name"`
	StakingIDNumber uint64 `json:"staking_id"`
}

// ParseCharity constructs Charity from an array of string
// array has to conform to data order in the pub78 txt
func ParseCharity(line []string) (*db.Charity, error) {
	if len(line) < 6 {
		return nil, errors.New("invalid record")
	}

	charity := db.Charity{
		StakingID:         line[0],
		Name:              line[1],
		City:              line[2],
		State:             line[3],
		Country:           line[4],
		DeductibilityCode: line[5],
		Is501c3:           true,
	}

	return &charity, nil
}

// NewCSVReader takes in a filename and returns a csv.Reader
func NewCSVReader(file io.Reader) *csv.Reader {
	reader := csv.NewReader(file)
	reader.Comma = '|'
	return reader
}
