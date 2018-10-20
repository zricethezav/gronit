package main

// gronit
// [go]cron[monitor]

import (
	"log"
	"os"
)

const EMPTYSTR string = ""

const AWS_KEY = "AKIALALEMEL33243OLIAE"
const AWS_SECRET = "99432bfewaf823ec3294e231"

var opts *Options

func main() {
	args := os.Args[1:]

	db, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	opts = parseOptions(args)
	serverStart(opts, db)
}
