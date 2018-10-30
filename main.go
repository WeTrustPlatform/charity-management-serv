package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
)

type Charity struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
	EIN     string `json:"ein,omitempty"`
	Website string `json:"website,omitempty"`
}

var charities []Charity

func GetCharities(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(charities)
}

func GetCharity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range charities {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Charity{})
}

// Load .env variables
func init() {
	gotenv.Load()
}

func main() {
	// TODO remove hardcoded data
	charities = append(charities, Charity{
		ID:      "1",
		Name:    "Open Library",
		Address: "Internet Archive, 300 Funston Avenue, San Francisco, CA 94118",
		EIN:     "94-3242767",
		Website: "https://openlibrary.org",
	})

	charities = append(charities, Charity{
		ID:      "2",
		Name:    "Code for America",
		Address: "155 9th Street, San Francisco, CA 94103 USA",
		EIN:     "27-1067272",
		Website: "https://www.codeforamerica.org",
	})

	router := mux.NewRouter()
	router.HandleFunc("/charities", GetCharities).Methods("GET")
	router.HandleFunc("/charities/{id}", GetCharity).Methods("GET")
	port := getEnv("PORT", "8001")
	log.Println("Listening on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
