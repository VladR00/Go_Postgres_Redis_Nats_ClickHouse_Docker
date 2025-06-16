package handlers

import (
	"encoding/json"
	"fmt"
	response "gopostgres/internal/domain/models/handle"
	postgres "gopostgres/pkg/storage/requestStorage"
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
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: fmt.Sprintf("Decode URL path. Want /good/create/{int}: %v", err)})
		return
	}
	//usecase -> менеджер транзакций (begin/commit/rollback)
	answer, err := postgres.NewStoragePostgres(s.Db).Create(request, id)
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
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: "Only PATCH method allowed"})
		return
	}

	var request response.UpdatePayload

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: "Decode. Want 'name':'string', 'description':'string'. Second - optional."})
		return
	}

	if request.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: "Decode. Want 'name':'string'"})
		return
	}

	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/good/update/"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: fmt.Sprintf("Decode URL path. Want /good/update/{int}: %v", err)})
		return
	}

	answer, err := postgres.NewStoragePostgres(s.Db).Update(request, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		log.Println(err)
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: fmt.Sprintf("Bad request: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answer)
}
func (s *StorageHandler) HandlerRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: "Only DELETE method allowed"})
		return
	}
	idstr := strings.TrimPrefix(r.URL.Path, "/good/remove/")
	idstrsplit := strings.Split(idstr, "&")
	if len(idstrsplit) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: `Decode URL path. Want /good/remove/{int}&{int}`})
		return
	}
	id, err := strconv.Atoi(idstrsplit[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: fmt.Sprintf(`Decode URL path. Want /good/remove/{int}&{int}: %v`, err)})
		return
	}

	projectid, err := strconv.Atoi(idstrsplit[1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: fmt.Sprintf(`Decode URL path. Want /good/remove/{int}&{int}: %v`, err)})
		return
	}
	answer, err := postgres.NewStoragePostgres(s.Db).Remove(id, projectid)
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
func (s *StorageHandler) HandlerList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: "Only GET method allowed"})
		return
	}

	idstr := strings.TrimPrefix(r.URL.Path, "/goods/list/")
	idstrsplit := strings.Split(idstr, "&")
	var limit, offset int
	if len(idstrsplit) < 2 {
		limit, err := strconv.Atoi(idstr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: fmt.Sprintf(`Decode URL path. Want /goods/list/{int}&{int}: %v`, err)})
			return
		}
		offset = 1
		answer, err := postgres.NewStoragePostgres(s.Db).List(limit, offset)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: fmt.Sprintf("Bad request: %v", err)})
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(answer)
		return
	}
	limit, err := strconv.Atoi(idstrsplit[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: fmt.Sprintf(`Decode URL path. Want /goods/list/{int}&{int}: %v`, err)})
		return
	}

	offset, err = strconv.Atoi(idstrsplit[1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: fmt.Sprintf(`Decode URL path. Want /goods/list/{int}&{int}: %v`, err)})
		return
	}
	answer, err := postgres.NewStoragePostgres(s.Db).List(limit, offset)
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
func (s *StorageHandler) HandlerReprioritize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: "Only PATCH method allowed"})
		return
	}

	idstr := strings.TrimPrefix(r.URL.Path, "/good/reprioritize/")
	idstrsplit := strings.Split(idstr, "&")
	if len(idstrsplit) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: `Decode URL path. Want /good/reprioritize/{int}&{int}`})
		return
	}

	var request response.ReoprioritizePayload

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: "Decode. Want 'name':'string', 'description':'string'. Second - optional."})
		return
	}

	if request.NewPriority == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.DefaultResponse{Type: "Error", Message: "Decode. Want 'name':'string'"})
		return
	}
	priority := *request.NewPriority
	//need to continue
}
