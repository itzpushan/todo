package db

import (
	"context"

	"github.com/itzpushan/todo/ent"
	"github.com/itzpushan/todo/internal/config"
	_ "github.com/lib/pq"
)

func NewClient() (*ent.Client, error) {
	dsn := config.GetEnv("DATABASE_URL", "postgresql://postgres:pushan@localhost:5432/postgres?sslmode=disable")

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, err
	}

	return client, nil
}
