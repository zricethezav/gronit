package main

import (
	_ "bufio"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	_ "syscall"
)

type Task struct {
	Name    string `yaml:"name" json:"name"`
	Second  string `yaml:"second" json:"second"`
	Minute  string `yaml:"minute" json:"minute"`
	Hour    string `yaml:"hour" json:"hour"`
	Day     string `yaml:"day" json:"day"`
	Month   string `yaml:"month" json:"month"`
	Command string `yaml:"command" json:"command"`
}

// getTasks parses the system infromation and options to build a slice
// of Tasks. A slice of Tasks is returned.
func getTasks(sys *System, opts *Options) []Task {
	var tasks []Task

	if opts.LoadYaml != EMPTYSTR {
		yamlToTasks(&tasks, opts)
	} else if opts.LoadJson != EMPTYSTR {
		// open yaml file
		// parse yaml for jobs
		jsonToTasks(&tasks, opts)
	} else if opts.LoadCron != EMPTYSTR {
		// open yaml file
		// parse yaml for jobs
	}
	return tasks
}

// yamlToTasks is a helper for getTasks. This function reads a yaml file
// then marshals it into an slice of Tasks.
func yamlToTasks(tasks *[]Task, opts *Options) {
	yamlFile, err := ioutil.ReadFile(opts.LoadYaml)
	if err != nil {
		log.Printf("Error reading %s: %v ", opts.LoadYaml, err)
	}
	err = yaml.Unmarshal(yamlFile, &tasks)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

// jsonToTasks is a helper for getTasks. This function reads a json file
// then marshals it into an slice of Tasks.
func jsonToTasks(tasks *[]Task, opts *Options) {
	jsonFile, err := ioutil.ReadFile(opts.LoadJson)
	if err != nil {
		log.Printf("Error reading %s: %v ", opts.LoadJson, err)
	}
	if err := json.Unmarshal(jsonFile, &tasks); err != nil {
		panic(err)
	}
}

// cronToTasks is a helper for getTasks. This function reads a cron text file
// then marshals it into a slice of Tasks. Note: example of a cron text is `crontab -l`
func cronToTask() *Task {
	// TODO
	return &Task{}
}

// taskToCron writes out tasks to crontab
func tasksToCron(tasks []Task, sys *System) {
	var cronBytes []byte
	var cronTaskString string

	pathToCron := fmt.Sprintf("%s%s", sys.CronPrefix, sys.User)
	content, err := ioutil.ReadFile(pathToCron)
	if err != nil {
		fmt.Printf("\nno crontab for %s. You may need to run `sudo crontab -e -u %s`\nto initiate cron\n\n",
			sys.User, sys.User)
		log.Fatal(err)
	}
	file, err := os.Open(pathToCron)
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	for _, task := range tasks {
		cronTaskString = fmt.Sprintf("%s %s %s %s %s %s\n",
			task.Second, task.Minute, task.Hour, task.Day,
			task.Month, task.Command)
		cronBytes = append(cronBytes, []byte(cronTaskString)...)
	}

	fmt.Printf("\nAdding the following cronjobs to crontab for %s:\n\n%s", sys.User, cronBytes[:])
	if err := ioutil.WriteFile(pathToCron,
		append(content[:], cronBytes[:]...), 0777); err != nil {
		log.Fatal(err)
	}
}
