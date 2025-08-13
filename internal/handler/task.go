package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryszhio/tasktracker/internal/model"
	"github.com/ryszhio/tasktracker/internal/repository"
	"github.com/ryszhio/tasktracker/internal/service"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req model.RequestCreateTask
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.taskService.CreateTask(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	t, err := h.taskService.GetTasks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	tid := c.Param("id")

	t, err := h.taskService.GetTaskByID(c.Request.Context(), tid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	uid := c.Query("uid")
	email := c.Query("email")
	if uid == "" && email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please pass a valid task id or an email."})
		return
	}

	var (
		t   []repository.Task
		err error
	)

	if uid != "" {
		t, err = h.taskService.GetTasksByUserID(c.Request.Context(), uid)
	} else {
		t, err = h.taskService.GetTasksByEmail(c.Request.Context(), email)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

func (h *TaskHandler) UpdateTaskDetails(c *gin.Context) {
	var req model.RequestUpdateTask
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tid := c.Param("id")
	res, err := h.taskService.UpdateTaskDetails(c.Request.Context(), tid, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please pass a valid task id as param."})
		return
	}

	if err := h.taskService.DeleteTask(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Deleted"})
}
