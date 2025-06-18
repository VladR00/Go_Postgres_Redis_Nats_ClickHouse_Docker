package storage

import (
	"context"
	"fmt"
	config "gopostgres/internal/config"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Db *pgxpool.Pool
}

func NewStorage(database *pgxpool.Pool) *Storage {
	return &Storage{Db: database}
}

// migrate -path ./internal/domain/migrations/ -database "postgres://user:pass@localhost:5042/database?sslmode=disable" up
func migrations(url string) error {
	m, err := migrate.New("file://../internal/domain/migrations/postgresql", url)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	log.Println("Migrations applied successfully!")
	return nil
}

func ConnectPostgreSQL() (*pgxpool.Pool, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Panic(err)
	}
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:5042/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.DbName)

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = cfg.PoolSize

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	if err := pingDB(pool); err != nil {
		return nil, err
	}
	if err := migrations(dbURL); err != nil {
		return nil, err
	}
	return pool, nil
}

func pingDB(pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return err
	}
	log.Println("Posgtresql connected")
	return nil
}
