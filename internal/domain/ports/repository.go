package ports

import (
	"task-management-backend/internal/domain/entity"
	"task-management-backend/pkg/constant"
)

type TaskRepository interface {
	GetAllByUserID(userID int64) ([]entity.Task, error)
	GetSubTasks(parentID int64) ([]entity.Task, error)
	GetByID(id, userID int64) (*entity.Task, error)
	Create(task *entity.Task) error
	Update(task *entity.Task) error
	Delete(id, userID int64) error
	GetByUserIDAndStatus(userID int64, status constant.TaskStatus) ([]entity.Task, error)
}

type UserRepository interface {
	GetByUsername(username string) (*entity.User, error)
	Create(user *entity.User) error
	Upsert(user *entity.User) (*entity.User, error)
}
