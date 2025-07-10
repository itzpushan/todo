package todos

import (
	"github.com/gorilla/mux"
	"github.com/itzpushan/todo/ent"
	"github.com/itzpushan/todo/internal/middleware"
)

func Register(r *mux.Router, client *ent.Client) {
	protected := r.NewRoute().Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/", CreateTodo(client)).Methods("POST")
	protected.HandleFunc("/", GetAllTodos(client)).Methods("GET")
	protected.HandleFunc("/{id}", GetTodoByID(client)).Methods("GET")
	protected.HandleFunc("/{id}", UpdateTodo(client)).Methods("PUT")
	protected.HandleFunc("/{id}", DeleteTodo(client)).Methods("DELETE")
}
