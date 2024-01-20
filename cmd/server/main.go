package main

import (
	"fmt"
	"log"
	"myapp/database"
	"myapp/internal/handler"
	"net/http"
)

func main() {
	db, err := database.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	database.Migrate()

	http.HandleFunc("/api/enrich", handler.HandleRequest)
	fmt.Println("Server is running on :8080...")
	http.ListenAndServe(":8080", nil)
}
