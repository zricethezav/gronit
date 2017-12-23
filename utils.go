package main

import (
	_ "bufio"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	_ "os"
	"os/exec"
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
func tasksToCron(tasks []Task, sys *System, opts *Options) {
	var (
		cronTaskString string
		crontab        []byte
		err            error
	)

	if crontab, err = exec.Command("crontab", "-l").Output(); err != nil {
		log.Fatal(err)
	}

	for _, task := range tasks {
		cronTaskString = fmt.Sprintf("%s %s %s %s %s %s\n",
			task.Second, task.Minute, task.Hour, task.Day,
			task.Month, task.Command)
		crontab = append(crontab, []byte(cronTaskString)...)
	}

	err = ioutil.WriteFile("/tmp/gronit", crontab, 0777)
	if err != nil {
		log.Fatal(err)
	}

	if sys.User == "root" && opts.User != EMPTYSTR {
		err = exec.Command("crontab", "-u", opts.User, "/tmp/gronit").Run()
	} else {
		err = exec.Command("crontab", "/tmp/gronit").Run()
	}
	if err != nil {
		fmt.Printf("user %s does not exist\n", sys.User)
		log.Fatal(err)
	}

	err = exec.Command("rm", "/tmp/gronit").Run()
	if err != nil {
		log.Fatal(err)
	}
}
