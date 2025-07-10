package todos

import (
	"encoding/json"
	"net/http"

	"github.com/itzpushan/todo/ent"
	"github.com/itzpushan/todo/ent/todo"
	"github.com/itzpushan/todo/ent/user"
	"github.com/itzpushan/todo/internal/middleware"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CreateTodo(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := middleware.GetUserID(r.Context())
		userID, _ := uuid.Parse(userIDStr)

		var req TodoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		todo, err := client.Todo.Create().
			SetID(uuid.New()).
			SetTitle(req.Title).
			SetDescription(req.Description).
			SetAuthorID(userID).
			Save(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(todo)
	}
}

func GetAllTodos(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := middleware.GetUserID(r.Context())
		userID, _ := uuid.Parse(userIDStr)

		todos, err := client.Todo.Query().
			Where(todo.HasAuthorWith(user.IDEQ(userID))).
			All(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(todos)
	}
}

func GetTodoByID(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.Parse(mux.Vars(r)["id"])
		userIDStr := middleware.GetUserID(r.Context())
		userID, _ := uuid.Parse(userIDStr)

		todo, err := client.Todo.Get(r.Context(), id)

		if err != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}

		if todo.AuthorID != userID {
			http.Error(w, "Unauthorized access", http.StatusForbidden)
			return
		}

		json.NewEncoder(w).Encode(todo)
	}
}

func UpdateTodo(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.Parse(mux.Vars(r)["id"])
		userIDStr := middleware.GetUserID(r.Context())
		userID, _ := uuid.Parse(userIDStr)

		originalTodo, err := client.Todo.Get(r.Context(), id)
		if err != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}

		if originalTodo.AuthorID != userID {
			http.Error(w, "Unauthorized access", http.StatusForbidden)
			return
		}

		var req TodoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		todo, err := client.Todo.UpdateOneID(id).
			SetTitle(req.Title).
			SetDescription(req.Description).
			Save(r.Context())

		if err != nil {
			http.Error(w, "Update failed", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(todo)
	}
}

func DeleteTodo(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.Parse(mux.Vars(r)["id"])
		userIDStr := middleware.GetUserID(r.Context())
		userID, _ := uuid.Parse(userIDStr)

		originalTodo, e := client.Todo.Get(r.Context(), id)
		if e != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}

		if originalTodo.AuthorID != userID {
			http.Error(w, "Unauthorized access", http.StatusForbidden)
			return
		}

		err := client.Todo.DeleteOneID(id).Exec(r.Context())

		if err != nil {
			http.Error(w, "Delete failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
