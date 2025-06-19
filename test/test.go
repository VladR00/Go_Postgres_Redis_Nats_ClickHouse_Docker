package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	_ "github.com/ClickHouse/clickhouse-go" // Импорт драйвера ClickHouse
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/nats-io/nats.go"
)

// Структура для лог-сообщения
type LogMessage struct {
	Id          uint32    `json:"Id"`
	ProjectId   uint32    `json:"ProjectId"`
	Name        string    `json:"Name"`
	Description string    `json:"Description"`
	Priority    uint32    `json:"Priority"`
	Removed     uint8     `json:"Removed"`
	EventTime   time.Time `json:"EventTime"`
}

func migrations(url string) error {
	log.Println("start migrations")
	m, err := migrate.New("file://internal/domain/migrations/clickhouse", url)
	if err != nil {
		return err
	}
	log.Println("try up")
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Println("error up")
		return err
	}
	log.Println("Migrations applied successfully!")
	return nil
}

func main() {
	// Подключение к NATS
	log.Println("start")
	natsConn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer natsConn.Close()
	log.Println("nats connected")

	// Подключение к ClickHouse
	db, err := sql.Open("clickhouse", "tcp://localhost:9000?username=user&password=12333&database=test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("click connected")

	if err := migrations("clickhouse://localhost:9000?username=user&password=12333&database=test"); err != nil {
		log.Println("fatal migrations")
		log.Fatal(err)
	}

	// Подписка на NATS
	log.Println("subscribe")
	natsConn.Subscribe("logs", func(m *nats.Msg) {
		var logMsg LogMessage
		err := json.Unmarshal(m.Data, &logMsg)
		if err != nil {
			log.Printf("Failed to unmarshal log message: %v", err)
			return
		}

		log.Printf("Received log message: %+v", logMsg)

		// Вставка в ClickHouse
		log.Println("click insert")
		tx, err := db.Begin() // Начинаем транзакцию
		if err != nil {
			log.Printf("Failed to begin transaction: %v", err)
			return
		}

		insert, err := db.Prepare("INSERT INTO logs (Id, ProjectId, Name, Description, Priority, Removed, EventTime) VALUES (?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Printf("Failed to prepare statement: %v", err)
			tx.Rollback() // Откатим транзакцию в случае ошибки
			return
		}
		defer insert.Close()

		_, err = insert.Exec(logMsg.Id, logMsg.ProjectId, logMsg.Name, logMsg.Description, logMsg.Priority, logMsg.Removed, logMsg.EventTime)
		if err != nil {
			log.Printf("Failed to insert log into ClickHouse: %v", err)
			tx.Rollback() // Откатим транзакцию в случае ошибки
			return
		}

		err = tx.Commit() // Фиксируем транзакцию
		if err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			return
		}

		log.Printf("Inserted log into ClickHouse: %+v", logMsg)
	})

	// Пример отправки сообщения в NATS
	logMsg := LogMessage{
		Id:          3,
		ProjectId:   1,
		Name:        "в",
		Description: "f",
		Priority:    3,
		Removed:     0,
		EventTime:   time.Now(),
	}

	data, err := json.Marshal(logMsg)
	if err != nil {
		log.Fatalf("Failed to marshal log message: %v", err)
	}

	err = natsConn.Publish("logs", data)
	if err != nil {
		log.Fatalf("Failed to publish message to NATS: %v", err)
	}

	log.Println("Sent message to NATS:", string(data))

	// Ожидание
	log.Println("Listening for log messages on NATS...")
	select {}
}
