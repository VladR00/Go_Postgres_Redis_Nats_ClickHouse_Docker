package main

import (
	"fmt"
	handlers "gopostgres/internal/handlers"
	storagepostgresql "gopostgres/pkg/storage"
	storageclickhouse "gopostgres/pkg/storage/clickhouse"
	natspackage "gopostgres/pkg/storage/nats"
	"log"
	"net/http"
)

func main() {
	postgresql, err := storagepostgresql.ConnectPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}
	defer postgresql.Close()

	clickhouse, err := storageclickhouse.ConnectClickHouse()
	if err != nil {
		log.Fatal(err)
	}
	defer clickhouse.Close()

	nats, err := natspackage.ConnectNats()
	if err != nil {
		log.Fatal(err)
	}

	natspackage.NatsSubscribes(nats, clickhouse)

	storageHanlder := handlers.NewStorageHandler(postgresql, nats)

	http.HandleFunc("/good/create/", storageHanlder.HandlerCreate)             // POST	curl -X POST http://localhost:8080/good/create/1 -H "Content-Type: application/json" -d '{"name":"name"}' | jq
	http.HandleFunc("/good/update/", storageHanlder.HandlerPatch)              // PATCH
	http.HandleFunc("/good/remove/", storageHanlder.HandlerRemove)             // DELETE
	http.HandleFunc("/goods/list/", storageHanlder.HandlerList)                // GET	curl -X GET http://localhost:8080/goods/list/100 | jq
	http.HandleFunc("/good/reprioritize/", storageHanlder.HandlerReprioritize) // PATCH curl -X PATCH http://localhost:8080/good/reprioritize/1\&1 -H "Content-Type: application/json" -d '{"newPriority":1}' | jq

	fmt.Println("Server start at 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
