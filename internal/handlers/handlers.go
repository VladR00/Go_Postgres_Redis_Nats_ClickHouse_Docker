package handlers

import (
	"encoding/json"
	"net/http"
	//storage "gopostgres/storage"
)

func HandlerCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DefaultResponse{"Error", "Only POST method allowed"})
		return
	}

	//var data storage.Data

	// if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	response = map[string]string{"error": "Bad request"}
	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(response)
	// 	return
	// }

	//data.AddQuote()
	response := DefaultResponse{"Message", "Successfully"}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func HandlerPatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DefaultResponse{"Error", "Only PATCH method allowed"})
		return
	}
}
func HandlerRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DefaultResponse{"Error", "Only DELETE method allowed"})
		return
	}
}
func HandlerList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DefaultResponse{"Error", "Only GET method allowed"})
		return
	}
}
func HandlerReprioritize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DefaultResponse{"Error", "Only PATCH method allowed"})
		return
	}
}
