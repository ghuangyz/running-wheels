package types

import (
	"fmt"
	"testing"
)

var filename = "../../testdata/task_test.yml"

func TestTaskLoading(t *testing.T) {
	taskTable, err := LoadTaskTable(filename)
	if err != nil {
		fmt.Println(ErrorStackTrace(err))
	} else {
		for _, task := range taskTable {
			fmt.Printf("Name: %s\n", task.Name)
			for _, cmd := range task.CommandGroup {
				fmt.Printf("%s %v\n", cmd.Name, cmd.Arguments)
			}
			fmt.Printf("Depends: %v\n", task.Depends)
			fmt.Printf("-------------------------------------------\n")
		}
	}

	fmt.Printf("\n")
}
