package cache

import (
	"fmt"
	"sync"
	"task-management-backend/internal/domain/entity"
	"task-management-backend/pkg/constant"
	"time"
)

type Item struct {
	Tasks      []entity.Task
	Expiration time.Time
}

type TaskCache struct {
	mu    sync.RWMutex
	tasks map[string]Item
	ttl   time.Duration
}

func NewTaskCache(ttl time.Duration) *TaskCache {
	cache := &TaskCache{
		tasks: make(map[string]Item),
		ttl:   ttl,
	}

	go cache.cleanupExpired()
	return cache
}

func generateKey(userID int64, status constant.TaskStatus) string {
	return fmt.Sprintf("%d:%s", userID, status)
}

func (c *TaskCache) Get(userID int64, status constant.TaskStatus) ([]entity.Task, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := generateKey(userID, status)
	tasks, ok := c.tasks[key]
	if !ok {
		return nil, false
	}

	if time.Now().After(tasks.Expiration) {
		return nil, false
	}

	return tasks.Tasks, true
}

func (c *TaskCache) Set(userID int64, status constant.TaskStatus, tasks []entity.Task) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := generateKey(userID, status)
	c.tasks[key] = Item{
		Tasks:      tasks,
		Expiration: time.Now().Add(c.ttl),
	}
}

func (c *TaskCache) Invalidate(userID int64, statuses []constant.TaskStatus) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, status := range statuses {
		key := generateKey(userID, status)
		delete(c.tasks, key)
	}
}

func (c *TaskCache) cleanupExpired() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, task := range c.tasks {
			if now.After(task.Expiration) {
				delete(c.tasks, key)
			}
		}

		c.mu.Unlock()
	}
}
