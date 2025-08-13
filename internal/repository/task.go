package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) GetTaskByID(c context.Context, taskID uuid.UUID) (*Task, error) {
	query := `
			SELECT id, name, description, status, created_at, updated_at
			FROM tasks
			WHERE id = $1
	`
	var task Task
	err := r.db.QueryRowContext(c, query, taskID).Scan(
		&task.ID,
		&task.Name,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.Description,
	)

	if err != nil {
		return nil, fmt.Errorf("get task by id: %v", err)
	}

	return &task, nil
}

func (r *TaskRepository) GetTasks(c context.Context) ([]Task, error) {
	query := `
			SELECT id, name, description, status, created_at, updated_at
			FROM tasks
	`

	var tasks []Task

	rows, err := r.db.QueryContext(c, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get tasks: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("get tasks: %v", err)
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get tasks: %v", err)
	}
	return tasks, nil
}

func (r *TaskRepository) GetTasksByUserID(c context.Context, userID uuid.UUID) ([]Task, error) {
	query := `
		SELECT id, name, description, status, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
	`
	var tasks []Task

	rows, err := r.db.QueryContext(c, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get tasks by user id: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("get tasks by user id: %v", err)
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get tasks by user id: %v", err)
	}

	return tasks, nil
}

func (r *TaskRepository) CreateTask(c context.Context, task *Task) (*Task, error) {
	query := `
			INSERT INTO tasks (name, description, status, user_id)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(c, query, task.Name, task.Description, task.Status, task.UserID).Scan(
		&task.ID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("insert task: %v", err)
	}

	return task, nil
}

func (r *TaskRepository) UpdateName(c context.Context, taskID uuid.UUID, name string) (*Task, error) {
	query := `
			UPDATE tasks SET name = $1, updated_at = NOW()
			WHERE
			id = $2
			RETURNING id, name, description, status, created_at, updated_at
	`
	var task Task
	err := r.db.QueryRowContext(c, query, name, taskID).Scan(
		&task.ID,
		&task.Name,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no task found")
		}

		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) UpdateDescription(c context.Context, taskID uuid.UUID, desc string) (*Task, error) {
	query := `
			UPDATE tasks SET description = $1, updated_at = NOW()
			WHERE
			id = $2
			RETURNING id, name, description, status, created_at, updated_at
	`

	var task Task
	err := r.db.QueryRowContext(c, query, desc, taskID).Scan(
		&task.ID,
		&task.Name,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no task found")
		}

		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) UpdateStatus(c context.Context, taskID uuid.UUID, status string) (*Task, error) {
	query := `
			UPDATE tasks SET status = $1, updated_at = NOW()
			WHERE
			id = $2
			RETURNING id, name, description, status, created_at, updated_at
	`

	var task Task
	err := r.db.QueryRowContext(c, query, status, taskID).Scan(
		&task.ID,
		&task.Name,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no task found")
		}

		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) DeleteTask(c context.Context, taskID uuid.UUID) error {
	result, err := r.db.ExecContext(c, "DELETE FROM tasks WHERE id = $1", taskID)
	if err != nil {
		return fmt.Errorf("delete user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows effected: %v", err)
	}

	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}
