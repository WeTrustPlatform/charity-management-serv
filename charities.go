package main

import ()

type Charity struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
	EIN     string `json:"ein,omitempty"`
	Website string `json:"website,omitempty"`
}

func GetCharities(w http.ResponseWriter, r *http.Request) {
	var charities []Charity
	if err := orm.GetAll(&charities, ""); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(charities)
}

func GetCharity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var charity Charity
	if err := orm.Get(&charity, "id = ?", params["id"]); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if charity == nil {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(&Charity{})
}
