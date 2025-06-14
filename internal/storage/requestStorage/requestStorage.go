package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	response "gopostgres/internal/models/handle"
)

type StoragePostgres struct {
	Db *pgxpool.Pool
}

func NewStoragePostgres(storage *pgxpool.Pool) *StoragePostgres {
	return &StoragePostgres{Db: storage}
}

func (s *StoragePostgres) CreateTable() {
	create1 := `CREATE TABLE IF NOT EXISTS 
		projects(
			id SERIAL PRIMARY KEY,
			name VARCHAR(255),
			created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`
	if _, err := s.Db.Exec(context.Background(), create1); err != nil {
		log.Fatal("Create projects error:", err)
	}
	log.Println("Table created: projects")
	create2 := `CREATE TABLE IF NOT EXISTS
		goods(
			id SERIAL PRIMARY KEY,
			project_id INTEGER,
			name VARCHAR(255),
			description VARCHAR(255),
			priority INTEGER,
			removed BOOL,
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

func (s *StoragePostgres) CreateGood(request response.CreatePayload) {

}
