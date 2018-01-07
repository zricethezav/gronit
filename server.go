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
	"regexp"
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
	http.HandleFunc("/", status) // default to list
	http.HandleFunc("/create", create)
	http.HandleFunc("/run/", run)
	http.HandleFunc("/complete/", complete)
	http.HandleFunc("/status", status)
	http.HandleFunc("/history", history)
	host := fmt.Sprintf("localhost:%s", strconv.Itoa(opts.Port))
	log.Fatal(http.ListenAndServe(host, nil))
}

func serverStop(sys *System, opts *Options) {
	// TODO find process server running on and stop
}

func serverRestart(sys *System, opts *Options) {
	// TODO find process server running on and restart
}

// create yooo
func create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		rand.Seed(time.Now().UTC().UnixNano())
		h := sha256.New()
		randomInt := rand.Intn(10000000)
		h.Write([]byte(fmt.Sprintf("%d", randomInt)))
		key := fmt.Sprintf("%x", h.Sum(nil)[:3])
		fmt.Fprintf(w, string(key))
		_ = initEntry(key, db)
	}
}

// getKey extracts a job key and returns an error if invalid url
func getKey(label string, r *http.Request) (string, error) {
	regex := fmt.Sprintf("^/(%s)/([a-zA-Z0-9]+)$", label)
	var validPath = regexp.MustCompile(regex)
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		return "", fmt.Errorf("could not extract key from path %s", r.URL.Path)
	}
	return m[2], nil
}

// status returns the status of the job
func status(w http.ResponseWriter, r *http.Request) {
	key, err := getKey("status", r)
	if err != nil {
		fmt.Println("error")
	}
	getStatus(key, db)
}

// run some things
func run(w http.ResponseWriter, r *http.Request) {
	key, err := getKey("run", r)
	if err != nil {
		fmt.Println("error")
	}
	setStatus(key, "running", time.Now(), db)
}

// complete
func complete(w http.ResponseWriter, r *http.Request) {
	key, err := getKey("complete", r)
	if err != nil {
		fmt.Println("error")
	}
	setStatus(key, "complete", time.Now(), db)
}

// history returns the full history of the job
func history(w http.ResponseWriter, r *http.Request) {
	_, err := getKey("history", r)
	if err != nil {
		fmt.Println("error")
	}
}
