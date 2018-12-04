package db

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	p "github.com/Prabandham/paginator"
	"github.com/WeTrustPlatform/charity-management-serv/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Charity Model
// Contain data fields provided by Publication 78 Data
// https://www.irs.gov/charities-non-profits/tax-exempt-organization-search-bulk-data-downloads
type Charity struct {
	ID                uint64     `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt         time.Time  `json:"created_at,omitempty"`
	UpdatedAt         time.Time  `json:"updated_at,omitempty"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
	Name              string     `json:"name,omitempty" gorm:"index"`                // pub78
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
	query := r.URL.Query()
	page := query.Get("page")
	if len(page) == 0 {
		page = "1"
	}

	// if query has ein
	// then return the exact match of ein
	// and ignore the rest queries
	ein := query.Get("ein")
	if len(ein) > 0 {
		var charity Charity
		if dbInstance.Where("EIN = ?", ein).First(&charity).RecordNotFound() {
			http.NotFound(w, r)
			return
		}

		if err = json.NewEncoder(w).Encode(charity); err != nil {
			util.LogError(err)
		}

		return
	}

	// search by anything
	search := query.Get("search")
	var dataSource *gorm.DB
	if len(search) > 0 {
		searchValue := fmt.Sprintf("%s:*", search)
		dataSource = dbInstance.Where("tsv @@ to_tsquery(?)", searchValue)
	} else {
		dataSource = dbInstance
	}

	paginator := p.Paginator{
		DB:      dataSource,
		OrderBy: []string{"Name ASC"},
		Page:    page,
		PerPage: util.GetEnv("PER_PAGE", "10"), // Don't want clients to load all records
	}

	var charities []Charity
	results := paginator.Paginate(&charities)

	if err := json.NewEncoder(w).Encode(results); err != nil {
		util.LogError(err)
	}
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

	if err = json.NewEncoder(w).Encode(charity); err != nil {
		util.LogError(err)
	}
}
