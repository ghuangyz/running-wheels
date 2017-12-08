package task

import (
	"github.com/yhuang69/running-wheels/pkg/types"
)

type Runner interface {
	Reloadable() bool
	Loaded() bool
	LoadTasks(taskTable types.TaskTable) error
	Run(numThreads int) (RunSummary, error)
	TasksInfo() string
}
