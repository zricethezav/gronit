package main

// gronit
// [go]cron[monitor]

import (
	"log"
	"os"
)

const EMPTYSTR string = ""
const AWSKEY string = ""

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
