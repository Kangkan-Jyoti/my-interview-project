package repository

import (
	"context"

	todov1 "my-interview-project/gen/todo/v1"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoRepository struct {
	db *pgxpool.Pool
}

func NewTodoRepository(db *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(ctx context.Context, todo *todov1.Todo) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO todos (id, title, completed) VALUES ($1, $2, $3)`,
		todo.Id, todo.Title, todo.Completed,
	)
	return err
}

func (r *TodoRepository) Get(ctx context.Context, id string) (*todov1.Todo, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, title, completed FROM todos WHERE id=$1`, id,
	)

	var todo todov1.Todo
	err := row.Scan(&todo.Id, &todo.Title, &todo.Completed)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *TodoRepository) List(ctx context.Context) ([]*todov1.Todo, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, title, completed FROM todos`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*todov1.Todo
	for rows.Next() {
		var t todov1.Todo
		if err := rows.Scan(&t.Id, &t.Title, &t.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, &t)
	}

	return todos, nil
}

func (r *TodoRepository) Complete(ctx context.Context, id string) (*todov1.Todo, error) {
	_, err := r.db.Exec(ctx,
		`UPDATE todos SET completed=true WHERE id=$1`, id,
	)
	if err != nil {
		return nil, err
	}
	return r.Get(ctx, id)
}
