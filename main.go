package main

// gronit
// [go]cron[monitor]

import (
	_ "flag"
	"os"
)

const EMPTYSTR string = ""

var sys *System

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		help()
		os.Exit(1)
	}

	sys = defaultSys()
	opts := parseOptions(sys, args)
	if opts.Start {
		serverStart(sys, opts)
	} else if opts.Restart {
		// Try to restart
	} else if opts.Stop {
		// Try to stop
	}

	// non service, using util to add tasks
	tasks := getTasks(sys, opts)
	if opts.LoadYaml != EMPTYSTR || opts.LoadJson != EMPTYSTR || opts.LoadCron != EMPTYSTR {
		tasksToCron(tasks, sys)
	}
}
