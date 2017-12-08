package task

import (
	"bytes"
	"fmt"
	"time"
)

type Result struct {
	StdOut  bytes.Buffer
	StdErr  bytes.Buffer
	Errors  []error
	Status  string
	Elapsed time.Duration
}

type DAGRunSummary struct {
	elapsed time.Duration
	results []*Result
	names   []string
}

func NewDAGRunSummary(count int) *DAGRunSummary {
	status := new(DAGRunSummary)
	status.results = make([]*Result, count)
	status.names = make([]string, count)
	return status
}

func (status *DAGRunSummary) AddTaskResult(id int, name string, result interface{}) {
	tmp, _ := result.(*Result)
	status.results[id] = tmp
	status.names[id] = name
}

func (status *DAGRunSummary) SetElapsed(elapsed time.Duration) {
	status.elapsed = elapsed
}

func (status *DAGRunSummary) Elapsed() time.Duration {
	return status.elapsed
}

func (status *DAGRunSummary) TaskElapsed() map[string]time.Duration {
	ret := make(map[string]time.Duration)
	for i, result := range status.results {
		name := status.names[i]
		ret[name] = result.Elapsed
	}
	return ret
}

func (status *DAGRunSummary) TaskOutput() map[string]string {
	ret := make(map[string]string)
	for i, result := range status.results {
		name := status.names[i]
		var buffer bytes.Buffer
		buffer.WriteString(fmt.Sprintf("%s [%s]\n", name, result.Status))
		buffer.WriteString("STDOUT:\n")
		buffer.WriteString(result.StdOut.String())
		buffer.WriteString("STDERR:\n")
		buffer.WriteString(result.StdErr.String())
		buffer.WriteString("GOLANG Errors:\n")
		for i, err := range result.Errors {
			buffer.WriteString(fmt.Sprintf("Error %d: ", i))
			buffer.WriteString(err.Error())
			buffer.WriteString("\n")
		}

		buffer.WriteString("\n")
		ret[name] = buffer.String()
	}
	return ret
}

func (status *DAGRunSummary) TaskStatus() map[string]string {
	ret := make(map[string]string)
	for i, result := range status.results {
		name := status.names[i]
		ret[name] = result.Status
	}
	return ret
}

func (status *DAGRunSummary) Scalability() float32 {
	var accumulated int64 = 0
	for _, result := range status.results {
		accumulated += result.Elapsed.Nanoseconds()
	}
	total := status.elapsed.Nanoseconds()

	return float32(accumulated) / float32(total)
}
