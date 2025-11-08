package repository

import (
	"database/sql"
	"fmt"
	"task-management-backend/internal/domain/entity"
	"task-management-backend/pkg/constant"
	"time"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) GetAllByUserID(userID int64) ([]entity.Task, error) {
	query := `
		SELECT id, user_id, parent_id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE user_id = ? AND parent_id IS NULL
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}

	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		err := rows.Scan(&task.ID, &task.UserID, &task.ParentID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		subTasks, err := r.GetSubTasks(task.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get subtasks: %w", err)
		}

		task.SubTasks = subTasks
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskRepository) GetSubTasks(parentID int64) ([]entity.Task, error) {
	query := `
		SELECT id, user_id, parent_id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE parent_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query subtasks: %w", err)
	}

	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		err := rows.Scan(&task.ID, &task.UserID, &task.ParentID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subtask: %w", err)
		}

		subTasks, err := r.GetSubTasks(task.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get nested subtasks: %w", err)
		}

		task.SubTasks = subTasks
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskRepository) GetByID(id, userID int64) (*entity.Task, error) {
	query := `
		SELECT id, user_id, parent_id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE id = ? AND user_id = ?
	`

	var task entity.Task
	err := r.db.QueryRow(query, id, userID).Scan(
		&task.ID, &task.UserID, &task.ParentID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}

		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	subTasks, err := r.GetSubTasks(task.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subtasks: %w", err)
	}

	task.SubTasks = subTasks

	return &task, nil
}

func (r *TaskRepository) Create(task *entity.Task) error {
	query := `
		INSERT INTO tasks (user_id, parent_id, title, description, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now
	result, err := r.db.Exec(query, task.UserID, task.ParentID, task.Title, task.Description, task.Status, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	task.ID = id
	return nil
}

func (r *TaskRepository) Update(task *entity.Task) error {
	query := `
		UPDATE tasks
		SET title = ?, description = ?, status = ?, parent_id = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`
	task.UpdatedAt = time.Now()
	result, err := r.db.Exec(query, task.Title, task.Description, task.Status, task.ParentID, task.UpdatedAt, task.ID, task.UserID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *TaskRepository) Delete(id, userID int64) error {
	query := `DELETE FROM tasks WHERE id = ? AND user_id = ?`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *TaskRepository) GetByUserIDAndStatus(userID int64, status constant.TaskStatus) ([]entity.Task, error) {
	query := `
		SELECT id, user_id, parent_id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE user_id = ? AND status = ? AND parent_id IS NULL
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks by status: %w", err)
	}

	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		err := rows.Scan(&task.ID, &task.UserID, &task.ParentID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		subTasks, err := r.GetSubTasks(task.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get subtasks: %w", err)
		}

		task.SubTasks = subTasks
		tasks = append(tasks, task)
	}

	return tasks, nil
}
