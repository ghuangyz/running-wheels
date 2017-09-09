package main

import (
	"flag"
	"fmt"
	"github.com/yhuang69/running-wheels/pkg/task"
	"github.com/yhuang69/running-wheels/pkg/types"
	"os"
)

type RunCommandOptions struct {
	Threads  *int
	Filename *string
}

func main() {
	runCommand := flag.NewFlagSet("run", flag.ExitOnError)
	rcOptions := RunCommandOptions{}
	rcOptions.Threads = runCommand.Int("threads", 4, "-threads=N")
	rcOptions.Filename = runCommand.String(
		"filename",
		"",
		"-filename=path/to/task/description/file",
	)

	if len(os.Args) < 2 {
		fmt.Println("Usage: running-wheels [command] [-options]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		runCommand.Parse(os.Args[2:])
		if *rcOptions.Filename == "" {
			fmt.Println("Usage: running-wheels run -filename=path/to/file [-options]")
			os.Exit(1)
		}
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	var runner task.Runner
	runner = task.NewDAGRunner()
	taskTable, err := types.LoadTaskTable(*rcOptions.Filename)
	if err != nil {
		fmt.Println(types.ErrorStackTrace(err))
		return
	}

	runner.LoadTasks(taskTable)
	status, err := runner.Run(*rcOptions.Threads)
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
