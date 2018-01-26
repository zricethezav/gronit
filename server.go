package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/rand"
	"net/http"
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

       Cron Monitoring

`

var mu sync.Mutex
var db *bolt.DB

// serverStart starts the gronit server which routes a few
// paths: add, update, list, remove, logs
func serverStart(opts *Options, _db *bolt.DB) {
	db = _db
	fmt.Printf("%s", serverStartMsg)
	http.HandleFunc("/", create)
	http.HandleFunc("/create", create)
	http.HandleFunc("/run/", run)
	http.HandleFunc("/complete/", complete)
	http.HandleFunc("/clear/", clear)
	http.HandleFunc("/status/", status)
	http.HandleFunc("/history/", history)
	http.HandleFunc("/summary/", summary)
	host := fmt.Sprintf("localhost:%s", strconv.Itoa(opts.Port))
	log.Fatal(http.ListenAndServe(host, nil))
}

// create new job monitor
func create(w http.ResponseWriter, r *http.Request) {
	type idResponse struct {
		ID string `json:"id"`
	}
	if r.Method == "GET" {
		rand.Seed(time.Now().UTC().UnixNano())
		h := sha256.New()
		randomInt := rand.Intn(10000000)
		h.Write([]byte(fmt.Sprintf("%d", randomInt)))
		id := fmt.Sprintf("%x", h.Sum(nil)[:3])
		err := initEntry(id, db, time.Now())
		if err != nil {
			http.Error(w, "failed to create entry", http.StatusForbidden)
		}
		w.Header().Set("Content-Type", "application/json")
		idJSON, err := json.Marshal(idResponse{ID: id})
		if err != nil {
			http.Error(w, "failed to create entry", http.StatusForbidden)
		}
		w.Write(idJSON)
	}
}

// getID extracts a job id and returns an error if invalid url
func getID(label string, r *http.Request) (string, error) {
	regex := fmt.Sprintf("^/(%s)/([a-zA-Z0-9]+)$", label)
	var validPath = regexp.MustCompile(regex)
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		return "", fmt.Errorf("could not extract id from path %s", r.URL.Path)
	}
	return m[2], nil
}

// run some things
func run(w http.ResponseWriter, r *http.Request) {
	id, err := getID("run", r)
	if err != nil {
		fmt.Println("error")
	}
	setStatus(id, "running", time.Now(), db)
}

// complete
func complete(w http.ResponseWriter, r *http.Request) {
	id, err := getID("complete", r)
	if err != nil {
		fmt.Println("error")
	}
	setStatus(id, "complete", time.Now(), db)
}

// clear
func clear(w http.ResponseWriter, r *http.Request) {
	id, err := getID("clear", r)
	if err != nil {
		fmt.Println("error")
	}
	deleteBucket(id, db)
}

// status returns the status of the job
func status(w http.ResponseWriter, r *http.Request) {
	id, err := getID("status", r)
	status, err := getStatus(id, db)
	statusJSON, err := json.Marshal(status)
	if err != nil {
		fmt.Println("Error retrieving history")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(statusJSON)
}

// history returns the full history of the job
func history(w http.ResponseWriter, r *http.Request) {
	id, err := getID("history", r)
	history, err := getHistory(id, db)
	historyJSON, err := json.Marshal(history)
	if err != nil {
		fmt.Println("Error retrieving history")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(historyJSON)
}

// summary returns the summary of the job
func summary(w http.ResponseWriter, r *http.Request) {
	id, err := getID("summary", r)
	summary, err := getSummary(id, db)
	summaryJSON, err := json.Marshal(summary)
	if err != nil {
		fmt.Println("Error retrieving history")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(summaryJSON)
}
