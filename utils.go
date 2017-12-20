package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	_ "syscall"
)

const EMPTYSTR string = ""

const usage = `usage: gronit [options]

    -u --user		Select which user to update cron
    --loadyaml		Yaml file that contains a schedule
    --loadjson		JSON file that contains a schedule
    --loadcron		cron file that contains a schedule
    --path        	path to crontab:
	  * default cron path for osx: /var/at/tabs/$USER
	  * default cron path for linux: /etc/cron.d/$USER

    --list-json 	Sends crontabs to stdout in json format
    --list-yaml 	Sends crontabs to stdout in yaml format
    -l, --list    	Sends crontabs to stdout in human readable form

    -v --version	Version

`

type System struct {
	CronPrefix string
	OS         string
	User       string
}

type Options struct {
	User     string
	Port     int
	LoadYaml string
	LoadJson string
	LoadCron string
}

type Task struct {
	Name    string `yaml:"name"`
	Second  string `yaml:"second"`
	Minute  string `yaml:"minute"`
	Hour    string `yaml:"hour"`
	Day     string `yaml:"day"`
	Month   string `yaml:"month"`
	Command string `yaml:"command"`
}

// defaultSys fills a System struct with path to the crontab directory,
// default username, and type of system (macOS or Linux from `uname`)
func defaultSys() *System {
	var (
		cronPrefix string
		uname      []byte
		err        error
	)

	_user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	userName := _user.Username

	if uname, err = exec.Command("uname").Output(); err != nil {
		log.Fatal(err)
	}
	unameString := strings.TrimSpace(string(uname))

	if unameString == "Darwin" {
		// cronPrefix = "/var/at/tabs/"
		cronPrefix = "/usr/lib/cron/tabs/"
	} else if unameString == "Linux" {
		cronPrefix = "/etc/cron.d/"
	}

	return &System{
		CronPrefix: cronPrefix,
		OS:         unameString,
		User:       userName,
	}
}

// help prints the usage string and exits
func help() {
	os.Stderr.WriteString(usage)
	os.Exit(1)
}

// optionsNextInt is a parseOptions helper that returns the value (int) of an option
// if valid.
func optionsNextInt(args []string, i *int) int {
	if len(args) > *i+1 {
		*i++
	} else {
		help()
	}
	argInt, err := strconv.Atoi(args[*i])
	if err != nil {
		fmt.Printf("Invalid %s option: %s\n", args[*i-1], args[*i])
		help()
	}
	return argInt
}

// optionsNextString is a parseOptions helper that returns the value (string) of an option
// if valid.
func optionsNextString(args []string, i *int) string {
	if len(args) > *i+1 {
		*i++
	} else {
		help()
	}
	return args[*i]
}

// parseOptions
func parseOptions(defaultSys *System, args []string) *Options {
	opts := &Options{}
	loadFile := false

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-u", "--user":
			opts.User = optionsNextString(args, &i)
			defaultSys.User = opts.User
		case "-p", "--port":
			opts.Port = optionsNextInt(args, &i)
		case "--loadyaml":
			if !(loadFile) {
				opts.LoadYaml = optionsNextString(args, &i)
				loadFile = true
			}
		case "--loadjson":
			if !(loadFile) {
				opts.LoadYaml = optionsNextString(args, &i)
				loadFile = true
			}
		case "--loadcron":
			if !(loadFile) {
				opts.LoadYaml = optionsNextString(args, &i)
				loadFile = true
			}
		default:
			fmt.Printf("Uknown option %s\n\n", arg)
			help()
			return nil
		}
	}

	if opts.User == "" {
		opts.User = defaultSys.User
	}

	return opts
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
func jsonToTasks() *Task {
	return &Task{}
}

// cronToTasks is a helper for getTasks. This function reads a cron text file
// then marshals it into a slice of Tasks. Note: example of a cron text is `crontab -l`
func cronToTask() *Task {
	return &Task{}
}

// taskToCron TODO
func tasksToCron(tasks []Task, sys *System) {
	var cronBytes []byte
	var cronTaskString string

	tmpfile, err := ioutil.TempFile("", "gronit_tmp")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	for _, task := range tasks {
		cronTaskString = fmt.Sprintf("%s %s %s %s %s %s\n",
			task.Second, task.Minute, task.Hour, task.Day,
			task.Month, task.Command)
		cronBytes = append(cronBytes, []byte(cronTaskString)...)
	}
	fmt.Printf("%s", cronBytes)

	if _, err := tmpfile.Write(cronBytes); err != nil {
		log.Fatal(err)
	}

	// TODO read current cron, check for duplicates, print suggestions, cp
	pathToCron := fmt.Sprintf("%s%s", sys.CronPrefix, sys.User)
	content, err := ioutil.ReadFile(pathToCron)
	fmt.Println("yooo")
	fmt.Println(pathToCron)
	fmt.Printf("%s", content)
	fmt.Println(content)
}
