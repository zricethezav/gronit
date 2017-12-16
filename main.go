package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// gronit
// [go]cron[monitor]

// BASIC IDEA OF UPDATING:
// https://unix.stackexchange.com/questions/241118/how-to-programmatically-add-new-crontab-file-without-replacing-previous-one

func main() {

	sys := defaultSys()
	fmt.Println(sys)
	parseArgs()
	os.Exit(1)

	// Subcommands
	startCommand := flag.NewFlagSet("start", flag.ExitOnError)
	restartCommand := flag.NewFlagSet("restart", flag.ExitOnError)
	stopCommand := flag.NewFlagSet("stop", flag.ExitOnError)

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		fmt.Printf(usage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "update":
		fmt.Printf("Updating")
		update()
	case "start":
		startCommand.Parse(os.Args[2:])
	case "stop":
		stopCommand.Parse(os.Args[2:])
	case "restart":
		restartCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func update() {
	// PROTOTYPE:
	// okay now we can write to cron, need to pass in arguments like:
	// 	-- yaml,json file or text
	// need to check if cron exists already
	// Writing tmp file
	fakeCron := []byte("* * * * * touch $HOME/yooo\n")
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	fmt.Println(tmpfile.Name())
	if _, err := tmpfile.Write(fakeCron); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	// Reading file
	content, err := ioutil.ReadFile("/var/spool/cron/crontabs/root")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", append(content[:], fakeCron[:]...))
	if err := ioutil.WriteFile("/var/spool/cron/crontabs/root",
		append(content[:], fakeCron[:]...), 0777); err != nil {
		log.Fatal(err)
	}
}
