package main

import (
	"fmt"
	"log"
	"net/http"

	handlers "gopostgres/internal/handlers"
	storage "gopostgres/pkg/storage"
)

func main() {
	db, err := storage.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	storage.NewStorage(db).Initiate()

	storageHanlder := handlers.NewStorageHandler(db)

	http.HandleFunc("/good/create/", storageHanlder.HandlerCreate)            // POST
	http.HandleFunc("/good/update/", storageHanlder.HandlerPatch)             // PATCH
	http.HandleFunc("/good/remove", storageHanlder.HandlerRemove)             // DELETE
	http.HandleFunc("/goods/list", storageHanlder.HandlerList)                // GET
	http.HandleFunc("/good/reprioritize", storageHanlder.HandlerReprioritize) // PATCH

	fmt.Println("Server start at 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
