package types

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	White  = "NotStarted"
	Yellow = "DependencyFailed"
	Green  = "Finished"
	Red    = "Failed"
)

const (
	DuplicateTaskError = "DuplicateTaskError"
	FileReadError      = "FileReadError"
	YamlReadError      = "YamlReadError"
)

type Command struct {
	Name      string   `yaml:"cmd" json:"cmd"`
	Arguments []string `yaml:"args,omitempty" json"args,omitempty"`
}

type TaskData struct {
	Name         string     `yaml:"name" json:"name"`
	CommandGroup []*Command `yaml:"commands" json:"commands"`
	Depends      []string   `yaml:"depends,omitempty" json:"depends,omitempty"`
	UsePipe      bool       `yaml:"pipe,omitempty" json:"pipe,omitempty"`
}

type TaskDataList struct {
	Tasks []*TaskData `yaml:"tasks" json:"tasks"`
}

type Task struct {
	TaskData
	Status string
	Id     int
}

type TaskTable map[string]*Task

func LoadTaskTable(configFilename string) (TaskTable, error) {
	taskDataList := TaskDataList{}
	yamlFile, err := ioutil.ReadFile(configFilename)
	if err != nil {
		return nil, NewError(FileReadError, err.Error())
	}

	err = yaml.Unmarshal(yamlFile, &taskDataList)
	if err != nil {
		return nil, NewError(YamlReadError, err.Error())
	}

	taskTable := TaskTable{}
	for id, taskData := range taskDataList.Tasks {
		name := taskData.Name
		task := &Task{TaskData: *taskData, Status: White, Id: id}
		if _, exist := taskTable[name]; exist {
			return nil, NewError(DuplicateTaskError, name)
		}
		taskTable[name] = task
	}

	return taskTable, nil
}
