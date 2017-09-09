package task

import (
	"github.com/yhuang69/running-wheels/pkg/types"
	"time"
)

type Runner interface {
	Reloadable() bool
	Loaded() bool
	LoadTasks(taskTable types.TaskTable) error
	Run(numThreads int) (RunStatus, error)
	TasksInfo() string
}

type RunStatus interface {
	AddTaskResult(id int, name string, result interface{})
	SetElapsed(elapsed time.Duration)
	Elapsed() time.Duration
	TaskElapsed() map[string]time.Duration
	TaskOutput() map[string]string
	TaskStatus() map[string]string
	Scalability() float32
}
