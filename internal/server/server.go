package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/ryszhio/tasktracker/db"
	"github.com/ryszhio/tasktracker/internal/handler"
	"github.com/ryszhio/tasktracker/internal/repository"
	"github.com/ryszhio/tasktracker/internal/service"
)

type Server struct {
	ip   string
	port int
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	ip := os.Getenv("IP")
	srv := &Server{
		ip:   ip,
		port: port,
	}

	fmt.Printf("Initialized server with: %s:%d\n", srv.ip, srv.port)

	db, err := db.NewDatabase()
	if err != nil {
		msg := fmt.Sprintf("Error while creating database pool: %s", err.Error())
		panic(msg)
	}

	taskRepo := repository.NewTaskRepository(db)
	userRepo := repository.NewUserRepository(db)

	taskService := service.NewTaskService(taskRepo, userRepo)
	userService := service.NewUserService(userRepo)

	taskHandler := handler.NewTaskHandler(taskService)
	userHandler := handler.NewUserHandler(userService)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      srv.RegisterRoutes(taskHandler, userHandler),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return server
}
