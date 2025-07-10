package router

import (
	"github.com/gorilla/mux"
	"github.com/itzpushan/todo/ent"
	"github.com/itzpushan/todo/internal/todos"
	"github.com/itzpushan/todo/internal/user"
)

func New(client *ent.Client) *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	user.Register(api.PathPrefix("/user").Subrouter(), client)
	todos.Register(api.PathPrefix("/todos").Subrouter(), client)

	return r
}
