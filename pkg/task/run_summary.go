package task

import (
	"time"
)

type RunSummary interface {
	AddTaskResult(id int, name string, result interface{})
	SetElapsed(elapsed time.Duration)
	Elapsed() time.Duration
	TaskElapsed() map[string]time.Duration
	TaskOutput() map[string]string
	TaskStatus() map[string]string
	Scalability() float32
}
