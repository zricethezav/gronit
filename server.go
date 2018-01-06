package main

import (
	"crypto/sha256"
	_ "encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	_ "io/ioutil"
	"log"
	"math/rand"
	"net/http"
	_ "os/exec"
	"strconv"
	"sync"
	"time"
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

var mu sync.Mutex
var db *bolt.DB

// serverStart starts the gronit server which routes a few
// paths: add, update, list, remove, logs
func serverStart(sys *System, opts *Options, _db *bolt.DB) {
	db = _db
	fmt.Printf("%s", serverStartMsg)
	http.HandleFunc("/", list) // default to list
	http.HandleFunc("/create", create)
	http.HandleFunc("/start", start)
	http.HandleFunc("/complete", complete)
	http.HandleFunc("/list", list)
	host := fmt.Sprintf("localhost:%s", strconv.Itoa(opts.Port))
	log.Fatal(http.ListenAndServe(host, nil))
}

func serverStop(sys *System, opts *Options) {
	// TODO find process server running on and stop
}

func serverRestart(sys *System, opts *Options) {
	// TODO find process server running on and restart
}

func create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		rand.Seed(time.Now().UTC().UnixNano())
		h := sha256.New()
		randomInt := rand.Intn(10000000)
		h.Write([]byte(fmt.Sprintf("%d", randomInt)))
		token := fmt.Sprintf("%x\n", h.Sum(nil)[:3])
		fmt.Fprintf(w, string(token))
		// todo add to cache
		// map of token:cronjob data
		_ = set(token, []byte(""))
	}

	// running
	// completed
}

func list(w http.ResponseWriter, r *http.Request) {
	// break out
	_ = listBucket(db)

}

func start(w http.ResponseWriter, r *http.Request) {
	// TODO finish and put in handlers
}

func complete(w http.ResponseWriter, r *http.Request) {
	// TODO finish and put in handlers
}
