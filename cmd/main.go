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

	http.HandleFunc("/good/create/", storageHanlder.HandlerCreate)             // POST
	http.HandleFunc("/good/update/", storageHanlder.HandlerPatch)              // PATCH
	http.HandleFunc("/good/remove/", storageHanlder.HandlerRemove)             // DELETE
	http.HandleFunc("/goods/list/", storageHanlder.HandlerList)                // GET
	http.HandleFunc("/good/reprioritize/", storageHanlder.HandlerReprioritize) // PATCH

	fmt.Println("Server start at 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
