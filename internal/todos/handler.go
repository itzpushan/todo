package todos

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/itzpushan/todo/ent"
	"github.com/itzpushan/todo/ent/todo"
	"github.com/itzpushan/todo/ent/user"
	"github.com/itzpushan/todo/internal/middleware"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Global in-memory storage (protected by mutex)
var (
	recentlyViewed = make(map[string][]uuid.UUID)
	recentMu       sync.Mutex
)

const maxRecent = 10

type TodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TodoResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func updateRecentlyViewed(userID uuid.UUID, todoID uuid.UUID) {
	recentMu.Lock()
	defer recentMu.Unlock()

	list := recentlyViewed[userID.String()]

	for i, id := range list {
		if id == todoID {
			list = append(list[:i], list[i+1:]...)
			break
		}
	}

	list = append([]uuid.UUID{todoID}, list...)

	if len(list) > maxRecent {
		list = list[:maxRecent]
	}

	recentlyViewed[userID.String()] = list
}

func GetRecentlyViewed(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := middleware.GetUserID(r.Context())
		userID, _ := uuid.Parse(userIDStr)

		recentMu.Lock()
		ids := recentlyViewed[userID.String()]
		recentMu.Unlock()

		var todosList []TodoResponse

		for _, id := range ids {
			todo, err := client.Todo.Get(r.Context(), id)
			if err == nil {
				todosList = append(todosList, TodoResponse{
					ID:          todo.ID.String(),
					Title:       todo.Title,
					Description: todo.Description,
				})
			}
		}

		json.NewEncoder(w).Encode(todosList)

	}
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

		_, err := client.Todo.Create().
			SetID(uuid.New()).
			SetTitle(req.Title).
			SetDescription(req.Description).
			SetAuthorID(userID).
			Save(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "todo created",
		})
	}
}

func GetAllTodos(client *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := middleware.GetUserID(r.Context())
		userID, _ := uuid.Parse(userIDStr)

		const pageSize = 5

		page := 1
		if p := r.URL.Query().Get("page"); p != "" {
			if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
				page = parsed
			}
		}

		offset := (page - 1) * pageSize

		todos, err := client.Todo.Query().
			Where(todo.HasAuthorWith(user.IDEQ(userID))).
			Limit(pageSize).
			Offset(offset).
			All(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var response []TodoResponse
		for _, todo := range todos {
			response = append(response, TodoResponse{
				ID:          todo.ID.String(),
				Title:       todo.Title,
				Description: todo.Description,
			})
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"page":     page,
			"per_page": pageSize,
			"todos":    response,
		})
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

		updateRecentlyViewed(userID, id)

		response := TodoResponse{
			ID:          todo.ID.String(),
			Title:       todo.Title,
			Description: todo.Description,
		}

		json.NewEncoder(w).Encode(response)
	}
}

func UpdateTodo(client *ent.Client) http.HandlerFunc {
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

		var req TodoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err := client.Todo.UpdateOneID(id).
			SetTitle(req.Title).
			SetDescription(req.Description).
			Save(r.Context())

		if err != nil {
			http.Error(w, "Update failed", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "todo updates successfully",
		})
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
