package main1

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8" // Импорт Redis
)

type LogMessage struct {
	Id          int       `json:"id"`
	ProjectId   int       `json:"project_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	EventTime   time.Time `json:"event_time"`
}

func main() {
	// Создайте контекст
	ctx := context.Background()

	// Подключение к Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Адрес вашего Redis сервера
	})
	defer rdb.Close()

	log.Println("Connected to Redis")

	// Создание структуры
	logMessage := LogMessage{
		Id:          1,
		ProjectId:   42,
		Name:        "Sample Project",
		Description: "This is a sample description.",
		Priority:    1,
		Removed:     false,
		EventTime:   time.Now(),
	}

	// Сериализация структуры в JSON
	logData, err := json.Marshal(logMessage)
	if err != nil {
		log.Fatalf("Failed to marshal log message: %v", err)
	}

	// Сохранение в Redis под ключом "log:1"
	err = rdb.Set(ctx, "log:1", logData, 0).Err()
	if err != nil {
		log.Fatalf("Failed to set key in Redis: %v", err)
	}

	log.Println("Log message saved to Redis")

	// Получение данных из Redis
	result, err := rdb.Get(ctx, "log:1").Result()
	if err != nil {
		log.Fatalf("Failed to get key from Redis: %v", err)
	}

	// Десериализация JSON строки в структуру LogMessage
	var retrievedLogMsg LogMessage
	err = json.Unmarshal([]byte(result), &retrievedLogMsg)
	if err != nil {
		log.Printf("Failed to unmarshal log message: %v", err)
		return
	}

	// Вывод данных
	log.Printf("Retrieved log message: %+v", retrievedLogMsg)
}
