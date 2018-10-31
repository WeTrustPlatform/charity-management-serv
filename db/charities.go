package db

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Charity Model
// Contain data fields provided by Publication 78 Data
// https://www.irs.gov/charities-non-profits/tax-exempt-organization-search-bulk-data-downloads
type Charity struct {
	ID                uint64     `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt         time.Time  `json:"created_at,omitempty"`
	UpdatedAt         time.Time  `json:"updated_at,omitempty"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
	Name              string     `json:"name,omitempty"`                             // pub78
	City              string     `json:"city,omitempty"`                             // pub78
	State             string     `json:"state,omitempty"`                            // pub78
	Country           string     `json:"country,omitempty"`                          // pub78
	EIN               string     `json:"ein,omitempty" gorm:"not null;unique_index"` // pub78
	DeductibilityCode string     `json:"deductibility_code,omitempty"`               // pub78
	Website           string     `json:"website,omitempty"`                          // optional
	Address           string     `json:"address,omitempty"`                          // optional
	ContactInfo       string     `json:"contact_info,omitempty"`                     // optional
}

// GetCharities returns all charities in the http response
func GetCharities(w http.ResponseWriter, r *http.Request) {
	var charities []Charity
	if err := dbInstance.Find(&charities).Error; err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(charities)
}

// GetCharity returns one charity whose ID is in the request param
// Throw 404 if record is not found
func GetCharity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var charity Charity
	if dbInstance.Where("ID = ?", id).First(&charity).RecordNotFound() {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(charity)
}
