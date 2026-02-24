package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	todov1 "my-interview-project/gen/todo/v1"
	"my-interview-project/internal/service"

	"connectrpc.com/connect"
)

// TodoHandler exposes TodoService as REST/JSON API for UI clients.
type TodoHandler struct {
	svc *service.TodoService
}

func NewTodoHandler(svc *service.TodoService) *TodoHandler {
	return &TodoHandler{svc: svc}
}

type todoJSON struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func toJSON(t *todov1.Todo) todoJSON {
	if t == nil {
		return todoJSON{}
	}
	return todoJSON{ID: t.Id, Title: t.Title, Completed: t.Completed}
}

func (h *TodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/todos")
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		if path == "" {
			h.listTodos(ctx, w)
			return
		}
		if len(parts) == 1 && parts[0] != "" {
			h.getTodo(ctx, w, parts[0])
			return
		}
	case http.MethodPost:
		if path == "" {
			h.createTodo(ctx, w, r)
			return
		}

	case http.MethodPut:
		if path == "" {
			h.completeTodo(ctx, w, r)
			return
		}
	}

	http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
}

func (h *TodoHandler) listTodos(ctx context.Context, w http.ResponseWriter) {
	resp, err := h.svc.ListTodos(ctx, connect.NewRequest(&todov1.ListTodosRequest{}))
	if err != nil {
		writeError(w, err)
		return
	}
	todos := make([]todoJSON, len(resp.Msg.Todos))
	for i, t := range resp.Msg.Todos {
		todos[i] = toJSON(t)
	}
	json.NewEncoder(w).Encode(todos)
}

func (h *TodoHandler) getTodo(ctx context.Context, w http.ResponseWriter, id string) {
	resp, err := h.svc.GetTodo(ctx, connect.NewRequest(&todov1.GetTodoRequest{Id: id}))
	if err != nil {
		writeError(w, err)
		return
	}
	json.NewEncoder(w).Encode(toJSON(resp.Msg.Todo))
}

func (h *TodoHandler) createTodo(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if body.Title == "" {
		http.Error(w, `{"error":"title required"}`, http.StatusBadRequest)
		return
	}

	resp, err := h.svc.CreateTodo(ctx, connect.NewRequest(&todov1.CreateTodoRequest{Title: body.Title}))
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toJSON(resp.Msg.Todo))
}

func (h *TodoHandler) completeTodo(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if body.ID == "" {
		http.Error(w, `{"error":"id required"}`, http.StatusBadRequest)
		return
	}
	resp, err := h.svc.CompleteTodo(ctx, connect.NewRequest(&todov1.CompleteTodoRequest{Id: body.ID}))
	if err != nil {
		writeError(w, err)
		return
	}
	json.NewEncoder(w).Encode(toJSON(resp.Msg.Todo))
}

func writeError(w http.ResponseWriter, err error) {
	if connectErr, ok := err.(*connect.Error); ok {
		switch connectErr.Code() {
		case connect.CodeNotFound:
			http.Error(w, `{"error":"todo not found"}`, http.StatusNotFound)
			return
		}
	}
	http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
}
