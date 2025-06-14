package handlers

import (
	"encoding/json"
	response "gopostgres/internal/models/handle"
	postgres "gopostgres/internal/storage/requestStorage"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StorageHandler struct {
	Db *pgxpool.Pool
}

func NewStorageHandler(storage *pgxpool.Pool) *StorageHandler {
	return &StorageHandler{Db: storage}
}

func (s *StorageHandler) HandlerCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: "Only POST method allowed"})
		return
	}
	var request response.CreatePayload

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: "Decode. Want 'name':'string'"})
		return
	}

	postgres.NewStoragePostgres(s.Db).CreateGood(request)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: "Only PATCH method allowed"})
}
func (s *StorageHandler) HandlerPatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: "Only PATCH method allowed"})
		return
	}
}
func (s *StorageHandler) HandlerRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: "Only DELETE method allowed"})
		return
	}
}
func (s *StorageHandler) HandlerList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: "Only GET method allowed"})
		return
	}
}
func (s *StorageHandler) HandlerReprioritize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: "Only PATCH method allowed"})
		return
	}
}
