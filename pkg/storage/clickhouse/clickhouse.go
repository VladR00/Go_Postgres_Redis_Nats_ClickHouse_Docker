package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	response "gopostgres/internal/domain/models/handle"
	"log"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type StorageClickhouse struct {
	Db *sql.DB
}

func NewStorageClickhouse(db *sql.DB) *StorageClickhouse {
	return &StorageClickhouse{Db: db}
}

func ConnectClickHouse() (*sql.DB, error) {
	url := "localhost:9000?username=user&password=12333&database=test"
	db, err := sql.Open("clickhouse", fmt.Sprintf("tcp://%s", url))
	if err != nil {
		return nil, err
	}

	if err := migrations(fmt.Sprintf("clickhouse://%s", url)); err != nil {
		return nil, err
	}
	return db, nil
}

func migrations(url string) error {
	m, err := migrate.New("file://../internal/domain/migrations/clickhouse", url)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	log.Println("ClickHouse: Migrations applied.")
	return nil
}

func (s *StorageClickhouse) InsertLog(logs response.NatsForClick) error {
	tx, err := s.Db.Begin()
	if err != nil {
		log.Println("Error start transaction in click:", err)
		return err
	}
	insert := (`INSERT INTO logs 
				(Id, ProjectId, Name, Description, Priority, Removed, EventTime) 
				VALUES (?, ?, ?, ?, ?, ?, ?)`)

	//for _, logg := range logs {
	_, err = tx.ExecContext(context.Background(), insert, logs.Id, logs.ProjectId, logs.Name, logs.Description, logs.Priority, logs.Removed, logs.EventTime)
	if err != nil {
		log.Println("Failed to insert log into ClickHouse: ", err)
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		log.Println("Failed to commit transaction: ", err)
		return err
	}
	//}

	log.Println("Logs inserted into ClickHouse")
	return nil
}
