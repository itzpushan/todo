package main

import (
	"log"
	"net/http"

	"github.com/itzpushan/todo/internal/config"
	"github.com/itzpushan/todo/internal/db"
	"github.com/itzpushan/todo/internal/router"
)

func main() {
	config.LoadEnv()

	client, err := db.NewClient()
	if err != nil {
		log.Fatalf("Failed to Connect: %v", err)
	}

	defer client.Close()

	r := router.New(client)

	log.Println("Starting server at localhost:8000")

	log.Fatal(http.ListenAndServe(":8000", r))
}
