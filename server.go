package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

const serverStartMsg = `
   _____                 _ _   
  / ____|               (_) |  
 | |  __ _ __ ___  _ __  _| |_ 
 | | |_ | '__/ _ \| '_ \| | __|
 | |__| | | | (_) | | | | | |_ 
  \_____|_|  \___/|_| |_|_|\__|

     Cron Monitoring System

`

// serverStart starts the gronit server which routes a few
// paths: add, update, list, remove, logs
func serverStart(sys *System, opts *Options) {
	fmt.Printf("%s", serverStartMsg)
	http.HandleFunc("/", list) // default to list
	http.HandleFunc("/add", add)
	http.HandleFunc("/update", update)
	http.HandleFunc("/list", list)
	http.HandleFunc("/remove", remove)
	http.HandleFunc("/logs", logs)

	host := fmt.Sprintf("localhost:%s", strconv.Itoa(opts.Port))
	log.Fatal(http.ListenAndServe(host, nil))

}

func serverStop(sys *System, opts *Options) {
	// TODO
}

func serverRestart(sys *System, opts *Options) {
	// TODO
}

func add(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	// TODO
	mu.Unlock()
}

func update(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	// TODO
	mu.Unlock()
}

func list(w http.ResponseWriter, r *http.Request) {
	// TODO finish and put in handlers
}

func remove(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	// TODO
	mu.Unlock()
}

func logs(w http.ResponseWriter, r *http.Request) {
	// TODO finish and put in handlers
}
