package model

import "time"

type RequestCreateTask struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UserID      string `json:"user_id"`
}

type ResponseCreateTask struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RequestUpdateTask struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}
