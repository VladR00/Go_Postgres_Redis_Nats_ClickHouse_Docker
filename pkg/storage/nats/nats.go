package nats

import (
	"database/sql"
	"encoding/json"
	response "gopostgres/internal/domain/models/handle"
	"gopostgres/pkg/storage/clickhouse"
	"log"

	"github.com/nats-io/nats.go"
)

func ConnectNats() (*nats.Conn, error) {
	conn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NatsSubscribes(natsConn *nats.Conn, database *sql.DB) {
	natsConn.Subscribe("logs", func(m *nats.Msg) {
		var logs response.NatsForClick
		err := json.Unmarshal(m.Data, &logs)
		if err != nil {
			log.Println("Failed to unmarshal log message:", err)
		}
		log.Println("Nats recieve msg. try to give click")
		db := clickhouse.NewStorageClickhouse(database)
		if err := db.InsertLog(logs); err != nil {
			log.Println("error insert into click:", err)
		}
	})
}
