package postgres

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"

	response "gopostgres/internal/domain/models/handle"
)

type StoragePostgres struct {
	Db       *pgxpool.Pool
	NatsConn *nats.Conn
}

func NewStoragePostgres(storage *pgxpool.Pool, nats *nats.Conn) *StoragePostgres {
	return &StoragePostgres{Db: storage, NatsConn: nats}
}

func (s *StoragePostgres) Create(request response.CreatePayload, id int) (*response.Goods, error) { //create/good
	tx, err := s.Db.Begin(context.Background())
	if err != nil {
		log.Println("Transaction begin error")
		return nil, err
	}

	insertcreategoods := `	WITH S AS (SELECT * FROM goods FOR UPDATE)
							INSERT INTO goods(name, description, project_id, priority, removed)
							VALUES($1, $2, $3, (SELECT COALESCE(MAX(priority), 0) + 1 FROM goods), $4) RETURNING id, project_id, name, description, priority, removed, created_at`

	var answer response.Goods

	err = tx.QueryRow(context.Background(), insertcreategoods, request.Name, "", id, false).Scan(&answer.ID, &answer.ProjectID, &answer.Name, &answer.Description, &answer.Priority, &answer.Removed, &answer.CreatedAt)
	if err != nil {
		if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
			log.Println("Transaction rollback error:")
			return nil, err
		}
		log.Println("Transaction rollbacked.")
		return nil, err
	} else {
		if commitErr := tx.Commit(context.Background()); commitErr != nil {
			log.Println("Transaction commit error")
			return nil, err
		}
		log.Println("Transaction commited.")
	}

	var removed uint8
	if answer.Removed {
		removed = 1
	} else {
		removed = 0
	}

	natsdata := response.NatsForClick{
		Id:          uint32(answer.ID),
		ProjectId:   uint32(answer.ProjectID),
		Name:        answer.Name,
		Description: answer.Description,
		Priority:    uint32(answer.Priority),
		Removed:     removed,
		EventTime:   time.Now(),
	}

	natsanswer, err := json.Marshal(natsdata)
	if err != nil {
		log.Println("Failed to marshal log message: ", err)
	}

	err = s.NatsConn.Publish("logs", natsanswer)
	if err != nil {
		log.Println("Failed to publish message to NATS: ", err)
	}

	return &answer, nil
}

func (s *StoragePostgres) Update(request response.UpdatePayload, id int, projectid int) (*response.Goods, error) { //update/good PATCH
	tx, err := s.Db.Begin(context.Background())
	if err != nil {
		log.Println("Transaction begin error")
		return nil, err
	}
	updategoods := `WITH S AS (SELECT * FROM goods WHERE name=$1 FOR UPDATE)
					UPDATE goods SET description=$2, name=$1 WHERE id=$3 AND project_id=$4
					RETURNING id, project_id, name, description, priority, removed, created_at`

	var answer response.Goods

	err = tx.QueryRow(context.Background(), updategoods, request.Name, request.Description, id, projectid).Scan(&answer.ID, &answer.ProjectID, &answer.Name, &answer.Description, &answer.Priority, &answer.Removed, &answer.CreatedAt)
	if err != nil {
		if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
			log.Println("Transaction rollback error:", err)
			return nil, err
		}
		log.Println("Transaction rollbacked.")
		return nil, err
	} else {
		if commitErr := tx.Commit(context.Background()); commitErr != nil {
			log.Println("Transaction commit error:", err)
			return nil, err
		}
		log.Println("Transaction commited.")
	}

	var removed uint8
	if answer.Removed {
		removed = 1
	} else {
		removed = 0
	}

	natsdata := response.NatsForClick{
		Id:          uint32(answer.ID),
		ProjectId:   uint32(answer.ProjectID),
		Name:        answer.Name,
		Description: answer.Description,
		Priority:    uint32(answer.Priority),
		Removed:     removed,
		EventTime:   time.Now(),
	}

	natsanswer, err := json.Marshal(natsdata)
	if err != nil {
		log.Println("Failed to marshal log message: ", err)
	}

	err = s.NatsConn.Publish("logs", natsanswer)
	if err != nil {
		log.Println("Failed to publish message to NATS: ", err)
	}

	return &answer, nil
}

func (s *StoragePostgres) Remove(id int, projectid int) (*response.DeleteResponse, error) { //goods/remove DELETE
	tx, err := s.Db.Begin(context.Background())
	if err != nil {
		log.Println("Transaction begin error")
		return nil, err
	}
	removegoods := `WITH S AS (SELECT removed FROM goods WHERE id = $1 AND project_id = $2  FOR UPDATE)
					UPDATE goods SET removed=true WHERE id = $1 AND project_id = $2
					RETURNING id, project_id, removed`

	var answer response.DeleteResponse

	err = tx.QueryRow(context.Background(), removegoods, id, projectid).Scan(&answer.ID, &answer.CampaignID, &answer.Removed)
	if err != nil {
		if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
			log.Println("Transaction rollback error:", rollbackErr)
			return nil, rollbackErr
		}
		log.Println("Transaction rollbacked.")
		return nil, err
	} else {
		if commitErr := tx.Commit(context.Background()); commitErr != nil {
			log.Println("Transaction commit error:", commitErr)
			return nil, commitErr
		}
		log.Println("Transaction commited.")
	}

	var removed uint8
	if answer.Removed {
		removed = 1
	} else {
		removed = 0
	}

	natsdata := response.NatsForClick{
		Id:        uint32(answer.ID),
		Removed:   removed,
		EventTime: time.Now(),
	}

	natsanswer, err := json.Marshal(natsdata)
	if err != nil {
		log.Println("Failed to marshal log message: ", err)
	}

	err = s.NatsConn.Publish("logs", natsanswer)
	if err != nil {
		log.Println("Failed to publish message to NATS: ", err)
	}

	return &answer, nil
}

func (s *StoragePostgres) List(limit int, offset int) (*response.GetListResponse, error) { //goods/remove GET
	var answer response.GetListResponse
	answer.Meta.Total = 0
	answer.Meta.Limit = limit
	answer.Meta.Offset = offset

	removegoods := `SELECT id, project_id, name, description, priority, removed, created_at FROM goods
					WHERE id >= $1
					ORDER BY id
					LIMIT $2`

	rows, err := s.Db.Query(context.Background(), removegoods, offset, limit)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &answer, nil
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	var gods []response.Goods
	removed := 0
	count := 0

	for rows.Next() {
		var god response.Goods
		if err := rows.Scan(&god.ID, &god.ProjectID, &god.Name, &god.Description, &god.Priority, &god.Removed, &god.CreatedAt); err != nil {
			log.Println("Err scan rows:", err)
			return nil, err
		}
		if god.Removed {
			removed++
		}
		gods = append(gods, god)
		count++
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	answer.Meta.Total = count
	answer.Meta.Removed = removed
	answer.Goods = gods

	return &answer, nil
}

func (s *StoragePostgres) Reprioritize(id, projectid, priority int) (*response.ReoprioritizeResponse, error) { //good/reoprioritize PATCH
	var answer response.ReoprioritizeResponse
	tx, err := s.Db.Begin(context.Background())
	if err != nil {
		log.Println("Transaction begin error")
		return nil, err
	}
	reprioritize := `WITH D AS (SELECT id FROM goods WHERE id >= $1 AND project_id = $2 FOR UPDATE),
		S AS (SELECT id,	CASE 
								WHEN row_number() OVER (ORDER BY id) = 1 THEN $3
								ELSE $3 + row_number() OVER (ORDER BY id) - 1
							END AS new_priority
	FROM goods WHERE id >= $1 AND project_id = $2)

	UPDATE goods SET
	priority = S.new_priority 
	FROM S 
	WHERE goods.id = S.id AND project_id = $2
	RETURNING goods.id, goods.priority`

	rows, err := tx.Query(context.Background(), reprioritize, id, projectid, priority)
	if err != nil {
		if err == pgx.ErrNoRows {
			if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
				log.Println("Transaction rollback error:", rollbackErr)
				return nil, err
			}
			log.Println("Transaction rollbacked.")
		}
		return nil, err
	}
	defer rows.Close()

	var priorities []response.Priorities

	for rows.Next() {
		var priority response.Priorities
		if err := rows.Scan(&priority.ID, &priority.Priority); err != nil {
			log.Println("Err scan rows:", err)
			return nil, err
		}
		natsdata := response.NatsForClick{
			Id:        uint32(priority.ID),
			Priority:  uint32(priority.Priority),
			EventTime: time.Now(),
		}

		natsanswer, err := json.Marshal(natsdata)
		if err != nil {
			log.Println("Failed to marshal log message: ", err)
		}

		err = s.NatsConn.Publish("logs", natsanswer)
		if err != nil {
			log.Println("Failed to publish message to NATS: ", err)
		}
		priorities = append(priorities, priority)
	}

	if err := rows.Err(); err != nil {
		if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
			log.Println("Transaction rollback error:", rollbackErr)
			return nil, err
		}
		log.Println("Transaction rollbacked.")
		return nil, err
	}
	if commitErr := tx.Commit(context.Background()); commitErr != nil {
		log.Println("Transaction commit error:", commitErr)
		return nil, commitErr
	}
	log.Println("Transaction commited.")
	answer.Priorities = priorities

	return &answer, nil
}
