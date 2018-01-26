package main

// gronit
// [go]cron[monitor]

import (
	"log"
	"os"
)

const EMPTYSTR string = ""

var opts *Options

func main() {
	args := os.Args[1:]

	db, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	opts = parseOptions(args)
	if opts.StartServer {
		// server
		serverStart(opts, db)
	} else {
		// client
		runCmd(opts)
	}
}
