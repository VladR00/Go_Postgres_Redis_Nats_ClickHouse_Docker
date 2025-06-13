package handlers

import (
	"encoding/json"
	"net/http"
	//storage "gopostgres/storage"
)

func HandlerCreate(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"error": "Only POST method allowed"}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
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

	w.WriteHeader(http.StatusOK)
	response = map[string]string{"message": "Quote successfully added"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func HandlerPatch(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"error": "Only PATCH method allowed"}

	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
}
func HandlerRemove(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"error": "Only DELETE method allowed"}

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
}
func HandlerList(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"error": "Only GET method allowed"}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
}
func HandlerReprioritize(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"error": "Only PATCH method allowed"}

	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
}
