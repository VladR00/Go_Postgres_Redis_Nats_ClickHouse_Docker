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
	"github.com/nats-io/nats.go"
)

type StorageHandler struct {
	Db       *pgxpool.Pool
	NatsConn *nats.Conn
}

func NewStorageHandler(storage *pgxpool.Pool, nats *nats.Conn) *StorageHandler {
	return &StorageHandler{Db: storage, NatsConn: nats}
}

func (s *StorageHandler) HandlerCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.DefaultResponse{Type: "Error", Message: "Only POST method allowed"}.Response(w, http.StatusMethodNotAllowed)
		return
	}
	var request response.CreatePayload

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.DefaultResponse{Type: "Error", Message: "Decode. Want 'name':'string'"}.Response(w, http.StatusBadRequest)
		return
	}

	if request.Name == "" {
		response.DefaultResponse{Type: "Error", Message: "Decode. Want 'name':'string'"}.Response(w, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/good/create/"))
	if err != nil {
		response.DefaultResponse{Type: "Error", Message: fmt.Sprintf("Decode URL path. Want /good/create/{int}: %v", err)}.Response(w, http.StatusBadRequest)
		return
	}

	answer, err := postgres.NewStoragePostgres(s.Db, s.NatsConn).Create(request, id)
	if err != nil {
		response.DefaultResponse{Type: "Error", Message: fmt.Sprintf("Bad request: %v", err)}.Response(w, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answer)
}

func (s *StorageHandler) HandlerPatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		response.DefaultResponse{Type: "Error", Message: "Only PATCH method allowed"}.Response(w, http.StatusMethodNotAllowed)
		return
	}

	var request response.UpdatePayload

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.DefaultResponse{Type: "Error", Message: "Decode. Want 'name':'string', 'description':'string'. Second - optional."}.Response(w, http.StatusBadRequest)
		return
	}

	if request.Name == "" {
		response.DefaultResponse{Type: "Error", Message: "Decode. Want 'name':'string'"}.Response(w, http.StatusBadRequest)
		return
	}

	idstr := strings.TrimPrefix(r.URL.Path, "/good/update/")
	idstrsplit := strings.Split(idstr, "&")
	if len(idstrsplit) < 2 {
		response.DefaultResponse{Type: "Error", Message: `Decode URL path. Want /good/update/{int}&{int}`}.Response(w, http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idstrsplit[0])
	if err != nil {
		response.DefaultResponse{Type: "Error", Message: fmt.Sprintf(`Decode URL path. Want /good/update/{int}&{int}: %v`, err)}.Response(w, http.StatusBadRequest)
		return
	}

	projectid, err := strconv.Atoi(idstrsplit[1])
	if err != nil {
		response.DefaultResponse{Type: "Error", Message: fmt.Sprintf(`Decode URL path. Want /good/update/{int}&{int}: %v`, err)}.Response(w, http.StatusBadRequest)
		return
	}

	answer, err := postgres.NewStoragePostgres(s.Db, s.NatsConn).Update(request, id, projectid)
	if err != nil {
		response.DefaultResponse{Type: "Error", Message: fmt.Sprintf("Bad request: %v", err)}.Response(w, http.StatusBadRequest)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answer)
}
func (s *StorageHandler) HandlerRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response.DefaultResponse{Type: "Error", Message: "Only DELETE method allowed"}.Response(w, http.StatusMethodNotAllowed)
		return
	}
	idstr := strings.TrimPrefix(r.URL.Path, "/good/remove/")
	idstrsplit := strings.Split(idstr, "&")
	if len(idstrsplit) < 2 {
		response.DefaultResponse{Type: "Error", Message: `Decode URL path. Want /good/remove/{int}&{int}`}.Response(w, http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idstrsplit[0])
	if err != nil {
		response.DefaultResponse{Type: "Error", Message: fmt.Sprintf(`Decode URL path. Want /good/remove/{int}&{int}: %v`, err)}.Response(w, http.StatusBadRequest)
		return
	}

	projectid, err := strconv.Atoi(idstrsplit[1])
	if err != nil {
		response.DefaultResponse{Type: "Error", Message: fmt.Sprintf(`Decode URL path. Want /good/remove/{int}&{int}: %v`, err)}.Response(w, http.StatusBadRequest)
		return
	}
	answer, err := postgres.NewStoragePostgres(s.Db, s.NatsConn).Remove(id, projectid)
	if err != nil {
		response.DefaultResponse{Type: "Error", Message: fmt.Sprintf("Bad request: %v", err)}.Response(w, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answer)
}
func (s *StorageHandler) HandlerList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.DefaultResponse{Type: "Error", Message: "Only GET method allowed"}.Response(w, http.StatusMethodNotAllowed)
		return
	}

	idstr := strings.TrimPrefix(r.URL.Path, "/goods/list/")
	idstrsplit := strings.Split(idstr, "&")
	var limit, offset int
	if len(idstrsplit) < 2 {
		limit, err := strconv.Atoi(idstr)
		if err != nil {
			response.DefaultResponse{Type: "Error", Message: fmt.Sprintf(`Decode URL path. Want /goods/list/{int}&{int}: %v`, err)}.Response(w, http.StatusBadRequest)
			return
		}
		offset = 1
		answer, err := postgres.NewStoragePostgres(s.Db, s.NatsConn).List(limit, offset)
		if err != nil {
			response.DefaultResponse{Type: "Error", Message: fmt.Sprintf("Bad request: %v", err)}.Response(w, http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(answer)
		return
	}
	limit, err := strconv.Atoi(idstrsplit[0])
	if err != nil {
		response.DefaultResponse{Type: "Error", Message: fmt.Sprintf(`Decode URL path. Want /goods/list/{int}&{int}: %v`, err)}.Response(w, http.StatusBadRequest)
		return
	}

	offset, err = strconv.Atoi(idstrsplit[1])
	if err != nil {
		response.DefaultResponse{Type: "Error", Message: fmt.Sprintf(`Decode URL path. Want /goods/list/{int}&{int}: %v`, err)}.Response(w, http.StatusBadRequest)
		return
	}
	answer, err := postgres.NewStoragePostgres(s.Db, s.NatsConn).List(limit, offset)
	if err != nil {
		response.DefaultResponse{Type: "Error", Message: fmt.Sprintf("Bad request: %v", err)}.Response(w, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answer)
}
func (s *StorageHandler) HandlerReprioritize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		response.DefaultResponse{Type: "Error", Message: "Only PATCH method allowed"}.Response(w, http.StatusMethodNotAllowed)
		return
	}

	idstr := strings.TrimPrefix(r.URL.Path, "/good/reprioritize/")
	idstrsplit := strings.Split(idstr, "&")
	if len(idstrsplit) < 2 {
		response.DefaultResponse{Type: "Error", Message: `Decode URL path. Want /good/reprioritize/{int}&{int}`}.Response(w, http.StatusBadRequest)
		return
	}

	var request response.ReoprioritizePayload

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.DefaultResponse{Type: "Error", Message: "Decode. Want 'newPriority':{int}"}.Response(w, http.StatusBadRequest)
		return
	}

	if request.NewPriority == nil {
		response.DefaultResponse{Type: "Error", Message: "Decode. Want 'newPriority':{int}"}.Response(w, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idstrsplit[0])
	if err != nil {
		response.DefaultResponse{Type: "Error", Message: fmt.Sprintf(`Decode URL path. Want /good/reprioritize/{int}&{int}: %v`, err)}.Response(w, http.StatusBadRequest)
		return
	}

	projectid, err := strconv.Atoi(idstrsplit[1])
	if err != nil {
		response.DefaultResponse{Type: "Error", Message: fmt.Sprintf(`Decode URL path. Want /good/reprioritize/{int}&{int}: %v`, err)}.Response(w, http.StatusBadRequest)
		return
	}

	priority := *request.NewPriority

	answer, err := postgres.NewStoragePostgres(s.Db, s.NatsConn).Reprioritize(id, projectid, priority)
	if err != nil {
		response.DefaultResponse{Type: "Error", Message: fmt.Sprintf("Bad request: %v", err)}.Response(w, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answer)
}
