package entity

import (
	"task-management-backend/pkg/constant"
	"time"
)

type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Task struct {
	ID          int64               `json:"id" db:"id"`
	UserID      int64               `json:"user_id" db:"user_id"`
	ParentID    *int64              `json:"parent_id,omitempty" db:"parent_id"`
	Title       string              `json:"title" db:"title"`
	Description string              `json:"description" db:"description"`
	Status      constant.TaskStatus `json:"status" db:"status"`
	CreatedAt   time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at" db:"updated_at"`
	SubTasks    []Task              `json:"sub_tasks,omitempty" db:"-"`
}

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description,omitempty"`
	ParentID    *int64 `json:"parent_id,omitempty"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
	ParentID    *int64  `json:"parent_id,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	UserID int64  `json:"user_id"`
}
