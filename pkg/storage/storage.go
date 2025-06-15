package storage

import (
	"context"
	"fmt"
	config "gopostgres/internal/config"
	postgres "gopostgres/pkg/storage/requestStorage"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Db *pgxpool.Pool
}

func NewStorage(database *pgxpool.Pool) *Storage {
	return &Storage{Db: database}
}

func (s *Storage) Initiate() {
	postgres.NewStoragePostgres(s.Db).CreateTable()
	log.Println("Tables successfully initiated")
}

// type Storage intreface{
// 	Initiate()
// 	Create(ctx) //
// 	Delete()
// 	Update()
// 	Get()
// }

// ctxx
func ConnectDB() (*pgxpool.Pool, error) {
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
