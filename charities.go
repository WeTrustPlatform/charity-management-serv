package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Charity struct {
	ID        uint64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Name      string `json:"name,omitempty"`
	Address   string `json:"address,omitempty"`
	EIN       string `json:"ein,omitempty"`
	Website   string `json:"website,omitempty"`
}

func GetCharities(w http.ResponseWriter, r *http.Request) {
	var charities []Charity
	if err := db.Find(&charities).Error; err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(charities)
}

func GetCharity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var charity Charity
	if db.Where("id = ?", id).First(&charity).RecordNotFound() {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(&Charity{})
}
