package main

import (
	_ "bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	_ "syscall"
)

const usage = `usage: gronit [arguments [options]] [options]
	
Arguments:
    start 		Starts gronit server
    restart 		Restarts gronit server
    stop 		Stops gronit server
	
Options:
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
    -p --port 		Port to run server on

`

type System struct {
	CronPrefix string
	OS         string
	User       string
}

type Options struct {
	Start    bool
	Stop     bool
	Restart  bool
	User     string
	Port     int
	LoadYaml string
	LoadJson string
	LoadCron string
}

// defaultSys fills a System struct with path to the crontab directory,
// default username, and type of system (macOS or Linux from `uname`)
func defaultSys() *System {
	var (
		cronPrefix string
		err        error
	)

	_user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	userName := _user.Username

	return &System{
		CronPrefix: cronPrefix,
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
		case "start":
			continue
		case "restart":
			continue
		case "stop":
			continue
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
				opts.LoadJson = optionsNextString(args, &i)
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
