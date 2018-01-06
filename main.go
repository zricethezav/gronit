package main

// gronit
// [go]cron[monitor]

import (
	_ "flag"
	"log"
	"os"
)

const EMPTYSTR string = ""

var sys *System
var opts *Options

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		help()
		os.Exit(1)
	}

	db, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sys = defaultSys()
	opts = parseOptions(sys, args)
	serverStart(sys, opts, db)
}
