package main

import (
	"fmt"
	"github.com/yhuang69/running-wheels/pkg/task"
	"github.com/yhuang69/running-wheels/pkg/types"
	"os"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	rwApp = kingpin.New("running-wheels", "A thread DAG task runner")

	// the run subcommand
	run = rwApp.Command("run", "Run a DAG task")
	runFilename = run.Flag("filename", "path to the file of tasks in YAML format").Required().String()
	runThreads = run.Flag("threads", "number of threads to use, default is 4").Default("4").Int()
)

func main() {
	switch kingpin.MustParse(rwApp.Parse(os.Args[1:])) {
	case run.FullCommand():
		runTask(*runFilename, *runThreads)
	}
}

func runTask(filename string, threads int) {
	var runner task.Runner
	runner = task.NewDAGRunner()
	taskTable, err := types.LoadTaskTable(filename)
	if err != nil {
		fmt.Println(types.ErrorStackTrace(err))
		return
	}

	runner.LoadTasks(taskTable)
	status, err := runner.Run(threads)
	if err != nil {
		fmt.Println(types.ErrorStackTrace(err))
	} else {
		fmt.Println("\n\nRun Summary:")
		for _, output := range status.TaskOutput() {
			fmt.Println(output)
		}

		fmt.Printf("Total Elapsed = %s\n", status.Elapsed())
		for name, duration := range status.TaskElapsed() {
			fmt.Printf("Time Elapsed On %s = %s\n", name, duration)
		}
		fmt.Printf("Calculated Scalability = %.4f\n", status.Scalability())
	}
}
