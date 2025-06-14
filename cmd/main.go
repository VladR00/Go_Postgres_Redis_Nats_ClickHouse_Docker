package main

import (
	"fmt"
	"log"
	"net/http"

	handlers "gopostgres/internal/handlers"

	posgtres "gopostgres/storage"
)

func main() {
	db, err := posgtres.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	posgtres.DB = db

	http.HandleFunc("/good/create", handlers.HandlerCreate)             // POST
	http.HandleFunc("/good/update", handlers.HandlerPatch)              // PATCH
	http.HandleFunc("/good/remove", handlers.HandlerRemove)             // DELETE
	http.HandleFunc("/goods/list", handlers.HandlerList)                // GET
	http.HandleFunc("/good/reprioritize", handlers.HandlerReprioritize) // PATCH

	fmt.Println("Server start at 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
