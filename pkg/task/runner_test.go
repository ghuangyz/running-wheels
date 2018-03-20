package task

import (
	"fmt"
	"github.com/ghuangyz/running-wheels/pkg/types"
	"testing"
)

var filename = "../../testdata/runner_test.yml"

func createRunnerFromData(runner Runner) error {
	taskTable, err := types.LoadTaskTable(filename)
	if err != nil {
		return err
	}

	runner.LoadTasks(taskTable)
	return err
}

func TestRunnerCreate(t *testing.T) {
	var runner Runner
	runner = NewDAGRunner()
	err := createRunnerFromData(runner)
	if err != nil {
		fmt.Println(types.ErrorStackTrace(err))
	} else {
		fmt.Println(runner.TasksInfo())
	}
}

func TestRunnerRun(t *testing.T) {
	var runner Runner
	runner = NewDAGRunner()
	err := createRunnerFromData(runner)
	if err != nil {
		fmt.Println(types.ErrorStackTrace(err))
	} else {
		if runner.Loaded() {
			status, err := runner.Run(2)
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
	}
}
