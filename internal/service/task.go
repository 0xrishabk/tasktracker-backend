package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/ryszhio/tasktracker/internal/model"
	"github.com/ryszhio/tasktracker/internal/repository"
)

type TaskService struct {
	taskRepo *repository.TaskRepository
	userRepo *repository.UserRepository
	timeout  time.Duration
}

func NewTaskService(taskRepo *repository.TaskRepository, userRepo *repository.UserRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
		userRepo: userRepo,
		timeout:  time.Duration(2) * time.Second,
	}
}

func (s *TaskService) CreateTask(c context.Context, req model.RequestCreateTask) (*model.ResponseCreateTask, error) {
	c, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	log.Printf("TaskService.CreateTask - Starting task creation for: %s", req.Name)

	if req.Name == "" {
		log.Print("TaskService.CreateTask - No name was provided for the task.")
		return nil, fmt.Errorf("name field for task is required")
	}

	if req.Status == "" {
		req.Status = "TO_DO"
	}

	t := &repository.Task{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		UserID:      req.UserID,
	}

	task, err := s.taskRepo.CreateTask(c, t)
	if err != nil {
		log.Printf("TaskService.CreateTask - Database error: %v", err)
		return nil, fmt.Errorf("failed to create task: %v", err)
	}

	log.Printf("TaskService.CreateTask - Task creation was successful: %s", task.ID)

	return &model.ResponseCreateTask{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}

func (s *TaskService) GetTaskByID(c context.Context, taskID string) (*repository.Task, error) {
	c, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	log.Printf("TaskService.GetTaskByID - Starting attempt to fetch task by ID: %s", taskID)

	tid, err := uuid.Parse(taskID)
	if err != nil {
		log.Printf("TaskService.GetTaskByID - UUID parsing error: %v", err)
		return nil, err
	}

	t, err := s.taskRepo.GetTaskByID(c, tid)
	if err != nil {
		log.Printf("TaskService.GetTaskByID - Database error: %v", err)
		return nil, err
	}

	log.Printf("TaskService.GetTaskByID - Successfully fetched task by ID: %s", taskID)
	return t, nil
}

func (s *TaskService) GetTasksByUserID(c context.Context, userID string) ([]repository.Task, error) {
	c, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	log.Printf("TaskService.GetTasksByUserID - Starting attempt to fetch tasks by UserID: %s", userID)

	uid, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("TaskService.GetTasksByUserID - UUID parsing error: %v", err)
		return nil, err
	}

	t, err := s.taskRepo.GetTasksByUserID(c, uid)
	if err != nil {
		log.Printf("TaskService.GetTasksByUserID - Database error: %v", err)
		return nil, err
	}

	log.Printf("TaskService.GetTasksByUserID - Successfully fetched tasks by UserID: %s", userID)
	return t, nil
}

func (s *TaskService) GetTasks(c context.Context) ([]repository.Task, error) {
	c, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	log.Printf("TaskService.GetTasks - Starting attempt to fetch tasks.")

	t, err := s.taskRepo.GetTasks(c)
	if err != nil {
		log.Printf("TaskService.GetTasks - Database error: %v", err)
		return nil, err
	}

	log.Printf("TaskService.GetTasks - Successfully fetched tasks.")
	return t, nil
}

func (s *TaskService) GetTasksByEmail(c context.Context, email string) ([]repository.Task, error) {
	c, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	log.Printf("TaskService.GetTasksByEmail - Starting attempt to fetch tasks by email: %s", email)

	userID, err := s.userRepo.GetUserIDByEmail(c, email)
	if err != nil {
		log.Printf("TaskService.GetTasksByEmail - Database error: %v", err)
		return nil, err
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("TaskService.GetTasksByEmail - UUID parsing error (user_id): %v", err)
		return nil, err
	}

	t, err := s.taskRepo.GetTasksByUserID(c, uid)
	if err != nil {
		log.Printf("TaskService.GetTasksByEmail - Database error: %v", err)
		return nil, err
	}

	log.Printf("TaskService.GetTasksByEmail - Successfully fetched tasks by email: %s", email)
	return t, nil
}

func (s *TaskService) UpdateTaskDetails(c context.Context, taskID string, req *model.RequestUpdateTask) (*model.ResponseCreateTask, error) {
	if req == nil {
		return nil, errors.New("nothing to update")
	}

	if req.Name == nil && req.Description == nil && req.Status == nil {
		return nil, errors.New("nothing to update")
	}

	c, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	log.Printf("TaskService.UpdateTaskDetails - Starting attempt to update task records.")

	var err error
	tid, err := uuid.Parse(taskID)
	if err != nil {
		log.Printf("TaskService.UpdateTaskDetails - UUID parsing error: %v", err)
		return nil, err
	}

	var task *repository.Task

	if req.Name != nil {
		log.Printf("\tTaskService.UpdateTaskDetails - Updating name.")
		task, err = s.taskRepo.UpdateName(c, tid, *req.Name)
		if err != nil {
			log.Printf("> \tTaskService.UpdateTaskDetails - Database error: %v", err)
			return nil, err
		}
	}
	if req.Description != nil {
		log.Printf("\tTaskService.UpdateTaskDetails - Updating description.")
		task, err = s.taskRepo.UpdateDescription(c, tid, *req.Description)
		if err != nil {
			log.Printf("> \tTaskService.UpdateTaskDetails - Database error: %v", err)
			return nil, err
		}
	}
	if req.Status != nil {
		log.Printf("\tTaskService.UpdateTaskDetails - Updating status.")
		task, err = s.taskRepo.UpdateStatus(c, tid, *req.Status)
		if err != nil {
			log.Printf("> \tTaskService.UpdateTaskDetails - Database error: %v", err)
			return nil, err
		}
	}

	log.Printf("TaskService.UpdateTaskDetails - Succesffuly updated task records.")

	return &model.ResponseCreateTask{
		ID:          task.ID,
		Name:        task.Name,
		Status:      task.Status,
		Description: task.Description,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}

func (s *TaskService) DeleteTask(c context.Context, taskID string) error {
	c, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tid, err := uuid.Parse(taskID)
	if err != nil {
		return err
	}

	if err := s.taskRepo.DeleteTask(c, tid); err != nil {
		return err
	}
	return nil
}

/*
func (s *TaskService) UpdateName(c context.Context, taskID string, name string) (*model.ResponseCreateTask, error) {
	c, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tid, err := uuid.Parse(taskID)
	if err != nil {
		return nil, err
	}

	task, err := s.taskRepo.UpdateName(c, tid, name)
	if err != nil {
		return nil, err
	}

	return &model.ResponseCreateTask{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}

func (s *TaskService) UpdateDescription(c context.Context, taskID string, desc string) (*model.ResponseCreateTask, error) {
	c, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tid, err := uuid.Parse(taskID)
	if err != nil {
		return nil, err
	}

	task, err := s.taskRepo.UpdateDescription(c, tid, desc)
	if err != nil {
		return nil, err
	}

	return &model.ResponseCreateTask{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}

func (s *TaskService) UpdateStatus(c context.Context, taskID string, desc string) (*model.ResponseCreateTask, error) {
	c, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tid, err := uuid.Parse(taskID)
	if err != nil {
		return nil, err
	}

	task, err := s.taskRepo.UpdateStatus(c, tid, desc)
	if err != nil {
		return nil, err
	}

	return &model.ResponseCreateTask{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}
*/
