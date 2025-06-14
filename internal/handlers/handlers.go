package handlers

import (
	"encoding/json"
	"fmt"
	response "gopostgres/internal/models/handle"
	postgres "gopostgres/internal/storage/requestStorage"
	"log"
	"net/http"
	"strconv"
	"strings"

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

	if request.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: "Decode. Want 'name':'string'"})
		return
	}

	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/good/create/"))
	log.Println(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: fmt.Sprintf("Decode URL path. Want /good/create/{int}: %v", err)})
		return
	}

	answer, err := postgres.NewStoragePostgres(s.Db).UpdateOrCreate(request, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: fmt.Sprintf("Bad request: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answer)
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
