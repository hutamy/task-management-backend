package ports

import (
	"task-management-backend/internal/domain/entity"
	"task-management-backend/pkg/constant"
)

type TaskCache interface {
	Get(userID int64, status constant.TaskStatus) ([]entity.Task, bool)
	Set(userID int64, status constant.TaskStatus, tasks []entity.Task)
	Invalidate(userID int64, statuses []constant.TaskStatus)
}
