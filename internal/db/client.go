package db

import (
	"context"

	"github.com/itzpushan/todo/ent"
	_ "github.com/lib/pq"
)

func NewClient() (*ent.Client, error) {
	client, err := ent.Open("postgres", "host=localhost port=5432 user=postgres password=pushan dbname=postgres sslmode=disable")
	if err != nil {
		return nil, err
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, err
	}

	return client, nil
}
