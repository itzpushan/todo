package main

import (
	"context"
	"log"

	"github.com/itzpushan/todo/internal/db"
	"github.com/itzpushan/todo/internal/seed"
)

func main() {
	client, err := db.NewClient()
	if err != nil {
		log.Fatalf("❌ failed to connect to DB: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Optional: create schema if not already created
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("❌ failed creating schema: %v", err)
	}

	// Run the seed logic
	if err := seed.Run(ctx, client); err != nil {
		log.Fatalf("❌ seeding failed: %v", err)
	}
}
