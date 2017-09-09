package task

import (
	"bytes"
	"fmt"
	"github.com/yhuang69/running-wheels/pkg/graph"
	"github.com/yhuang69/running-wheels/pkg/types"
	"github.com/yhuang69/running-wheels/pkg/utils"
	"log"
	"time"
)

const (
	RunnerNotLoadedError = "RunnerNotLoadedError"
	CyclicTaskGraphError = "CyclicTaskGraphError"
	NoSuchTaskError      = "NoSuchTaskError"
)

///////////////////////////////////////////////////////////////////////////////
// DAGRunner
type DAGRunner struct {
	loaded             bool
	taskGraph          *graph.Graph
	degrees            map[*graph.Node]int
	workloads          chan *graph.Node
	feedbacks          chan *graph.Node
	producerCompletion bool
}

func NewDAGRunner() *DAGRunner {
	runner := &DAGRunner{loaded: false, producerCompletion: false}
	return runner
}

func (runner *DAGRunner) Reloadable() bool {
	return true
}

func (runner *DAGRunner) Loaded() bool {
	return runner.loaded
}

func (runner *DAGRunner) LoadTasks(taskTable types.TaskTable) error {
	tGraph, err := runner.createGraph(taskTable)
	if err != nil {
		return err
	}

	runner.taskGraph = tGraph
	runner.degrees = make(map[*graph.Node]int)
	for _, node := range tGraph.Nodes() {
		runner.degrees[node] = 0
	}

	for _, node := range tGraph.Nodes() {
		for _, neighbor := range node.Neighbors() {
			runner.degrees[neighbor]++
		}
	}

	runner.loaded = true
	return nil
}

func (runner *DAGRunner) Run(numThreads int) (RunStatus, error) {
	if !runner.loaded {
		return nil, types.NewError(RunnerNotLoadedError, "Runner can not Run()!")
	}

	for _, node := range runner.taskGraph.Nodes() {
		task, _ := node.Value.(*types.Task)
		task.Status = types.White
	}

	status := NewDAGRunStatus(runner.taskGraph.Size())
	startTime := time.Now()

	// MT - For feedbacks the general rule is to make it equal to the size of workers or larger
	//      to avoid deadlocking
	//    - The workloads channel should *ideally* be an unbounded channel, we make its bound to
	//      be the size of graph which essentially means unbounded
	runner.workloads = make(chan *graph.Node, runner.taskGraph.Size())
	runner.feedbacks = make(chan *graph.Node, numThreads)

	producerDone := make(chan bool, 1)
	workerDone := make(chan bool, numThreads)
	runner.runProducer(producerDone)
	runner.runWorkers(workerDone, numThreads, status)
	<-producerDone
	for i := 0; i < numThreads; i++ {
		<-workerDone
	}

	status.SetElapsed(time.Since(startTime))
	return status, nil
}

func (runner *DAGRunner) TasksInfo() string {
	var buffer bytes.Buffer
	for node, degree := range runner.degrees {
		task, _ := node.Value.(*types.Task)
		buffer.WriteString(fmt.Sprintf("%s : %d\n", task.Name, degree))
	}
	return buffer.String()
}

///////////////////////////////////////////////////////////////////////////////
// Private methods on DAGRunner

func (runner *DAGRunner) createGraph(taskTable types.TaskTable) (*graph.Graph, error) {
	taskGraph := new(graph.Graph)
	nodeTable := make(map[string]*graph.Node)
	for name, task := range taskTable {
		nodeTable[name] = taskGraph.MakeNode(task)
	}

	for name, task := range taskTable {
		to := nodeTable[name]
		for _, depend := range task.Depends {
			from, exist := nodeTable[depend]
			if !exist {
				return nil, types.NewError(NoSuchTaskError, depend)
			}
			err := taskGraph.AddEdge(from, to)
			if err != nil {
				return nil, err
			}
		}
	}

	if cyclic, path := taskGraph.HasCycle(); cyclic {
		return nil, types.NewError(CyclicTaskGraphError, pathString(path))
	}

	return taskGraph, nil
}

func (runner *DAGRunner) runProducer(done chan bool) {
	go func() {
		runner.seedTasks()
		for len(runner.degrees) > 0 {
			taskNode := <-runner.feedbacks
			task := taskNode.Value.(*types.Task)
			for _, neighbor := range taskNode.Neighbors() {
				if task.Status != types.Green {
					nTask, _ := neighbor.Value.(*types.Task)
					nTask.Status = types.Yellow
				}

				runner.degrees[neighbor]--
				if runner.degrees[neighbor] == 0 {
					runner.workloads <- neighbor
					delete(runner.degrees, neighbor)
				}
			}
		}

		// Cleanning up. Close the workloads channel,
		// signal main/caller thread about completion
		close(runner.workloads)
		done <- true
		log.Println("Producer completed!")
	}()
}

func (runner *DAGRunner) runWorkers(done chan bool, workerCount int, status *DAGRunStatus) {
	for i := 0; i < workerCount; i++ {
		go func(id int) {
			for taskNode := range runner.workloads {
				start := time.Now()
				var stdout bytes.Buffer
				var stderr bytes.Buffer
				var errors []error

				task, _ := taskNode.Value.(*types.Task)
				if task.Status == types.White {
					log.Printf("Worker %d [Running] %s\n", id, task.Name)
					errors = utils.RunCommandGroup(
						task.CommandGroup,
						&stdout,
						&stderr,
						task.UsePipe,
					)
					if len(errors) != 0 {
						task.Status = types.Red
					} else {
						task.Status = types.Green
					}
					log.Printf("Worker %d [Finished] %s [%s]\n", id, task.Name, task.Status)
				} else {
					log.Printf("Worker %d [Skipping] %s [%s]\n", id, task.Name, task.Status)
				}

				result := &Result{
					StdOut:  stdout,
					StdErr:  stderr,
					Errors:  errors,
					Status:  task.Status,
					Elapsed: time.Since(start),
				}
				status.AddTaskResult(task.Id, task.Name, result)
				if taskNode.HasNeighbor() {
					runner.feedbacks <- taskNode
				}
			}
			done <- true
			log.Printf("Worker %d completed!\n", id)
		}(i)
	}
}

func (runner *DAGRunner) seedTasks() {
	for node, degree := range runner.degrees {
		if degree == 0 {
			runner.workloads <- node
			delete(runner.degrees, node)
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
// private Helper functions

func pathString(path []*graph.Node) string {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	for i, node := range path {
		task, _ := node.Value.(*types.Task)
		buffer.WriteString(task.Name)
		if i+1 < len(path) {
			buffer.WriteString("->")
		}
	}
	buffer.WriteString("]")
	return buffer.String()
}
