package postgres

import (
	"context"
	"fmt"
	config "gopostgres/internal/config"
	"log"

	"github.com/jackc/pgx/v5"
)

func ConnectDB() (*pgx.Conn, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Panic(err)
	}
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:5042/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.DbName)

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}

	if err = pingDB(conn); err != nil {
		return nil, err
	}
	return conn, nil
}

func CloseDB(conn *pgx.Conn) {
	conn.Close(context.Background())
}

func pingDB(conn *pgx.Conn) error {
	if err := conn.Ping(context.Background()); err != nil {
		return err
	}
	return nil
}
