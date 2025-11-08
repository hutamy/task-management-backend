package constant

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "to do"
	TaskStatusInProgress TaskStatus = "in progress"
	TaskStatusDone       TaskStatus = "done"
	TaskStatusAll        TaskStatus = "all"
	TaskStatusDefault    TaskStatus = ""
)
