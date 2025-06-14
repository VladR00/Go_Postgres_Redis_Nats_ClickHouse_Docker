package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	response "gopostgres/internal/domain/models/handle"
)

type StoragePostgres struct {
	Db *pgxpool.Pool
}

func NewStoragePostgres(storage *pgxpool.Pool) *StoragePostgres {
	return &StoragePostgres{Db: storage}
}

func (s *StoragePostgres) CreateTable() { //migrations --
	create1 := `CREATE TABLE IF NOT EXISTS 
		projects(
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`
	if _, err := s.Db.Exec(context.Background(), create1); err != nil {
		log.Fatal("Create projects error:", err)
	}
	log.Println("Table created: projects")
	create2 := `CREATE TABLE IF NOT EXISTS
		goods(
			id SERIAL PRIMARY KEY,
			project_id INTEGER NOT NULL,
			name VARCHAR(255) NOT NULL,
			description VARCHAR(255),
			priority INTEGER NOT NULL,
			removed BOOL NOT NULL,
			created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) 
		)`
	if _, err := s.Db.Exec(context.Background(), create2); err != nil {
		log.Fatal("Create goods error:", err)
	}
	log.Println("Table created: goods")
	index1 := `CREATE INDEX IF NOT EXISTS idx_project_id ON goods(project_id)`
	if _, err := s.Db.Exec(context.Background(), index1); err != nil {
		log.Fatal("Create INDEX goods(project_id) error:", err)
	}
	log.Println("Index created: goods(project_id)")
	index2 := `CREATE INDEX IF NOT EXISTS idx_name ON goods(name)`
	if _, err := s.Db.Exec(context.Background(), index2); err != nil {
		log.Fatal("Create INDEX goods(name) error:", err)
	}
	log.Println("Indexes created: goods(name)")
	insert := `INSERT INTO projects (name)
		SELECT 'first'
		WHERE (SELECT COUNT(*) FROM projects) = 0
	`
	if _, err := s.Db.Exec(context.Background(), insert); err != nil {
		log.Fatal("Insert first error:", err)
	}
	log.Println("Insert into projects: first")
}

// передача TX(опционально), без response
func (s *StoragePostgres) Upsert(request response.CreatePayload, id int) (*response.CreateResponse, error) { //create/good
	tx, err := s.Db.Begin(context.Background())
	if err != nil {
		log.Println("Transaction begin error")
		return nil, err
	}

	//select for update // блокировка таблиц
	//take priority

	insertcreategoods := `INSERT INTO goods(name, description, project_id, priority, removed)
							VALUES($1, $2, $3, $4, $5)`

	_, err = tx.Exec(context.Background(), insertcreategoods, request.Name, "", id, 1, false)
	if err != nil {
		if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
			log.Println("Transaction rollback error:")
			log.Println(err)
		}
		log.Println("Transaction rollbacked.")
	} else {
		if commitErr := tx.Commit(context.Background()); commitErr != nil {
			log.Println("Transaction commit error")
			log.Println(err)
		}
		log.Println("Transaction commited.")
	}
	//без звезды
	selectgoods := `SELECT * FROM goods 
	WHERE name = $1 AND project_id = $2
	ORDER BY id DESC
	LIMIT 1;`

	var answer response.CreateResponse

	err = s.Db.QueryRow(context.Background(), selectgoods, request.Name, id).Scan(&answer.Goods.ID, &answer.Goods.ProjectID, &answer.Goods.Name, &answer.Goods.Description, &answer.Goods.Priority, &answer.Goods.Removed, &answer.Goods.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("Good not found.")
		} else {
			log.Println("Query failed: ", err)
		}
		return nil, err
	}
	return &answer, nil
}
