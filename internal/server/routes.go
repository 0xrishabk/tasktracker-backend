package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/0xrishabk/tasktracker/internal/handler"
)

func (s *Server) RegisterRoutes(taskHandler *handler.TaskHandler, userHandler *handler.UserHandler) http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	intializeUserRoutes(r, userHandler)
	initializeTaskRoutes(r, taskHandler)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Task manager api.",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "It's aight mate",
		})
	})

	return r
}

func intializeUserRoutes(r *gin.Engine, h *handler.UserHandler) {
	user := r.Group("/api/user")

	user.POST("/register", h.CreateUser)
	user.POST("/login", h.Login)
	user.DELETE("/:id", h.Delete)
}

func initializeTaskRoutes(r *gin.Engine, h *handler.TaskHandler) {
	task := r.Group("/api/task")

	task.POST("/", h.CreateTask)
	task.GET("/all-task", h.GetAllTasks)
	task.GET("/id/:id", h.GetTaskByID)
	task.GET("/user", h.GetTasks)
	task.DELETE("/:id", h.DeleteTask)
}
