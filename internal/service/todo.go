package service

import (
	"context"
	"fmt"
	"time"

	todov1 "my-interview-project/gen/todo/v1"
	"my-interview-project/internal/repository"

	"connectrpc.com/connect"
)

type TodoService struct {
	repo *repository.TodoRepository
}

func NewTodoService(repo *repository.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) CreateTodo(
	ctx context.Context,
	req *connect.Request[todov1.CreateTodoRequest],
) (*connect.Response[todov1.CreateTodoResponse], error) {

	todo := &todov1.Todo{
		Id:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Title:     req.Msg.Title,
		Completed: false,
	}

	if err := s.repo.Create(ctx, todo); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&todov1.CreateTodoResponse{
		Todo: todo,
	}), nil
}

func (s *TodoService) GetTodo(
	ctx context.Context,
	req *connect.Request[todov1.GetTodoRequest],
) (*connect.Response[todov1.GetTodoResponse], error) {

	todo, err := s.repo.Get(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&todov1.GetTodoResponse{
		Todo: todo,
	}), nil
}

func (s *TodoService) ListTodos(
	ctx context.Context,
	req *connect.Request[todov1.ListTodosRequest],
) (*connect.Response[todov1.ListTodosResponse], error) {

	todos, err := s.repo.List(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&todov1.ListTodosResponse{
		Todos: todos,
	}), nil
}

func (s *TodoService) CompleteTodo(
	ctx context.Context,
	req *connect.Request[todov1.CompleteTodoRequest],
) (*connect.Response[todov1.CompleteTodoResponse], error) {

	todo, err := s.repo.Complete(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&todov1.CompleteTodoResponse{
		Todo: todo,
	}), nil
}
