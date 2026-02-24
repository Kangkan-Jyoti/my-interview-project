package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"my-interview-project/gen/todo/v1/todov1connect"
	rest "my-interview-project/internal/api/handler"
	"my-interview-project/internal/api/middleware"
	"my-interview-project/internal/db"
	"my-interview-project/internal/repository"
	"my-interview-project/internal/service"
)

func main() {

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not set")
	}
	if u, err := url.Parse(dbURL); err == nil && u.Query().Get("sslmode") == "" {
		sep := "?"
		if strings.Contains(dbURL, "?") {
			sep = "&"
		}
		dbURL = dbURL + sep + "sslmode=disable"
	}

	ctx := context.Background()
	var dbpool *pgxpool.Pool
	for i := 0; i < 30; i++ {
		if err := db.RunMigrations(dbURL); err != nil {
			log.Printf("waiting for postgres: %v", err)
			time.Sleep(time.Second)
			continue
		}
		var err error
		dbpool, err = pgxpool.New(ctx, dbURL)
		if err != nil {
			log.Printf("waiting for postgres: %v", err)
			time.Sleep(time.Second)
			continue
		}
		break
	}
	if dbpool == nil {
		log.Fatal("failed to connect to postgres after 30 attempts")
	}
	defer dbpool.Close()

	repo := repository.NewTodoRepository(dbpool)
	svc := service.NewTodoService(repo)

	mux := http.NewServeMux()
	path, handler := todov1connect.NewTodoServiceHandler(svc)
	mux.Handle(path, handler)

	mux.Handle("/api/todos", middleware.CORS(rest.NewTodoHandler(svc)))
	mux.Handle("/api/todos/", middleware.CORS(rest.NewTodoHandler(svc)))

	log.Println("Todo service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
