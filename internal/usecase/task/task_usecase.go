package task

import (
	"fmt"
	"task-management-backend/internal/domain/entity"
	"task-management-backend/internal/domain/ports"
	"task-management-backend/pkg/constant"
)

type TaskUseCase struct {
	repo  ports.TaskRepository
	cache ports.TaskCache
}

func NewTaskUseCase(repo ports.TaskRepository, cache ports.TaskCache) *TaskUseCase {
	return &TaskUseCase{
		repo:  repo,
		cache: cache,
	}
}

func (uc *TaskUseCase) GetTasks(userID int64, status constant.TaskStatus) ([]entity.Task, error) {
	// check to cache first before query to database
	if cachedTasks, ok := uc.cache.Get(userID, status); ok {
		return cachedTasks, nil
	}

	var tasks []entity.Task
	var err error

	switch status {
	case constant.TaskStatusDefault, constant.TaskStatusAll:
		tasks, err = uc.repo.GetAllByUserID(userID)
	case constant.TaskStatusTodo, constant.TaskStatusInProgress, constant.TaskStatusDone:
		tasks, err = uc.repo.GetByUserIDAndStatus(userID, status)
	default:
		return nil, fmt.Errorf("invalid status filter: %s", status)
	}

	if err != nil {
		return nil, err
	}

	uc.cache.Set(userID, status, tasks)
	return tasks, nil
}

func (uc *TaskUseCase) CreateTask(userID int64, title, description string, parentID *int64) (*entity.Task, error) {
	if title == "" {
		return nil, fmt.Errorf("task title cannot be empty")
	}

	task := &entity.Task{
		UserID:      userID,
		ParentID:    parentID,
		Title:       title,
		Description: description,
		Status:      constant.TaskStatusTodo,
		SubTasks:    make([]entity.Task, 0),
	}

	if err := uc.repo.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	uc.cache.Invalidate(userID, []constant.TaskStatus{
		constant.TaskStatusTodo,
		constant.TaskStatusAll,
		constant.TaskStatusDefault,
	})
	return task, nil
}

func (uc *TaskUseCase) UpdateTask(userID, taskID int64, title, description *string, status *constant.TaskStatus, parentID *int64) (*entity.Task, error) {
	task, err := uc.repo.GetByID(taskID, userID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	oldStatus := task.Status

	if title != nil {
		if *title == "" {
			return nil, fmt.Errorf("task title cannot be empty")
		}

		task.Title = *title
	}

	if description != nil {
		task.Description = *description
	}

	if status != nil {
		task.Status = *status
	}

	if parentID != nil {
		// validate that task is not creating a circular relationship
		if *parentID != 0 {
			if err := uc.validateNoCircularRelationship(taskID, *parentID, userID); err != nil {
				return nil, err
			}
		}

		task.ParentID = parentID
	}

	if err := uc.repo.Update(task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	statusesToInvalidate := []constant.TaskStatus{
		oldStatus,
		task.Status,
		constant.TaskStatusAll,
		constant.TaskStatusDefault,
	}
	uc.cache.Invalidate(userID, statusesToInvalidate)
	return task, nil
}

func (uc *TaskUseCase) DeleteTask(userID, taskID int64) error {
	task, err := uc.repo.GetByID(taskID, userID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	if err := uc.repo.Delete(taskID, userID); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	uc.cache.Invalidate(userID, []constant.TaskStatus{
		task.Status,
		constant.TaskStatusAll,
		constant.TaskStatusDefault,
	})
	return nil
}

func (uc *TaskUseCase) GetTaskByID(userID, taskID int64) (*entity.Task, error) {
	task, err := uc.repo.GetByID(taskID, userID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	return task, nil
}

func (uc *TaskUseCase) validateNoCircularRelationship(taskID, newParentID, userID int64) error {
	// check if newParentID is the same as taskID
	if newParentID == taskID {
		return fmt.Errorf("a task cannot be its own parent")
	}

	allTasks, err := uc.repo.GetAllByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to validate relationship: %w", err)
	}

	taskMap := make(map[int64]*entity.Task)
	for i := range allTasks {
		taskMap[allTasks[i].ID] = &allTasks[i]
	}

	// check if newParentID is a descendant of taskID
	if uc.isDescendant(taskID, newParentID, taskMap) {
		return fmt.Errorf("cannot set parent: circular relationship detected")
	}

	return nil
}

func (uc *TaskUseCase) isDescendant(ancestorID, potentialDescendantID int64, taskMap map[int64]*entity.Task) bool {
	current := taskMap[potentialDescendantID]
	visited := make(map[int64]bool)

	for current != nil {
		if visited[current.ID] {
			return false
		}

		visited[current.ID] = true
		if current.ParentID != nil && *current.ParentID == ancestorID {
			return true
		}

		if current.ParentID != nil {
			current = taskMap[*current.ParentID]
		} else {
			current = nil
		}
	}

	for _, task := range taskMap {
		if task.ParentID != nil && *task.ParentID == ancestorID {
			if task.ID == potentialDescendantID {
				return true
			}

			if uc.isDescendantRecursive(task.ID, potentialDescendantID, taskMap, visited) {
				return true
			}
		}
	}

	return false
}

func (uc *TaskUseCase) isDescendantRecursive(parentID, targetID int64, taskMap map[int64]*entity.Task, visited map[int64]bool) bool {
	if visited[parentID] {
		return false
	}

	visited[parentID] = true
	for _, task := range taskMap {
		if task.ParentID != nil && *task.ParentID == parentID {
			if task.ID == targetID {
				return true
			}

			if uc.isDescendantRecursive(task.ID, targetID, taskMap, visited) {
				return true
			}
		}
	}

	return false
}
