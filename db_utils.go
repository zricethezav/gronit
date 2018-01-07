package main

import (
	"encoding/json"
	_ "flag"
	"fmt"
	"github.com/boltdb/bolt"
	_ "log"
	"time"
)

type Entry struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}

// initEntry
func initEntry(key string, db *bolt.DB) error {
	var err error
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(key))
		if err != nil {
			return fmt.Errorf("could not create root bucket: %v", err)
		}
		return nil
	})
	return err
}

// setStatus sets the state of the job
func setStatus(key string, statusStr string, date time.Time, db *bolt.DB) error {
	s := Entry{Status: statusStr, Time: date}
	statusBytes, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("could not marshal entry json: %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte(key)).Put([]byte("status"),
			statusBytes)
		if err != nil {
			return fmt.Errorf("could not update run: %v", err)
		}
		return nil
	})

	//set history
	history, err := getHistory(key, db)
	history = append(history, s)
	historyBytes, err := json.Marshal(history)
	err = db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte(key)).Put([]byte("history"),
			historyBytes)
		if err != nil {
			return fmt.Errorf("could not update run: %v", err)
		}
		return nil
	})

	return err
}

// getStatus grabs the status of the job
func getStatus(key string, db *bolt.DB) (*Entry, error) {
	status := Entry{}
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(key)).Get([]byte("status"))
		json.Unmarshal(b, &status)
		return nil
	})
	return &status, err
}

// getHistory grabs the status of the job
func getHistory(key string, db *bolt.DB) ([]Entry, error) {
	history := []Entry{}
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(key)).Get([]byte("history"))
		json.Unmarshal(b, &history)
		return nil
	})
	return history, err
}

// setupDB
func setupDB() (*bolt.DB, error) {
	db, err := bolt.Open("gronit.db", 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}
	fmt.Println("DB Setup Done")
	return db, nil
}
