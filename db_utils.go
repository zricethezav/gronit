package main

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"time"
)

type Entry struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}

// initEntry
func initEntry(id string, db *bolt.DB) error {
	var err error
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(id))
		if err != nil {
			return fmt.Errorf("could not bucket for %s: %v", id, err)
		}
		return nil
	})
	return err
}

// setStatus sets the state of the job
func setStatus(id string, statusStr string, date time.Time, db *bolt.DB) error {
	entry := Entry{Status: statusStr, Time: date}
	statusBytes, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal status into json: %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		jobBucket := tx.Bucket([]byte(id))
		if jobBucket == nil {
			return fmt.Errorf("failed to find record of %s", id)
		}
		err := jobBucket.Put([]byte("status"), statusBytes)
		if err != nil {
			return fmt.Errorf("failed to update status for %s: %v", id, err)
		}
		return nil
	})
	err = setHistory(id, &entry, db)
	return err

}

// setHistory updates the history of a job entry
func setHistory(id string, entry *Entry, db *bolt.DB) error {
	history, err := getHistory(id, db)
	history = append(history, *entry)
	historyBytes, err := json.Marshal(history)
	if err != nil {
		return fmt.Errorf("failed to marshal history into json: %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		jobBucket := tx.Bucket([]byte(id))
		if jobBucket == nil {
			return fmt.Errorf("failed to find record of %s", id)
		}
		err := jobBucket.Put([]byte("history"),
			historyBytes)
		if err != nil {
			return fmt.Errorf("failed to update history for  %s: %v", id, err)
		}
		return nil
	})
	return err
}

// getStatus grabs the status of the job
func getStatus(id string, db *bolt.DB) (*Entry, error) {
	status := Entry{}
	err := db.View(func(tx *bolt.Tx) error {
		jobBucket := tx.Bucket([]byte(id))
		if jobBucket == nil {
			return fmt.Errorf("failed to find record of %s", id)
		}
		statusBytes := jobBucket.Get([]byte("status"))
		json.Unmarshal(statusBytes, &status)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &status, nil
}

// getHistory grabs the status of the job
func getHistory(id string, db *bolt.DB) ([]Entry, error) {
	history := []Entry{}
	err := db.View(func(tx *bolt.Tx) error {
		jobBucket := tx.Bucket([]byte(id))
		if jobBucket == nil {
			return fmt.Errorf("failed to find record of %s", id)
		}
		b := jobBucket.Get([]byte("history"))
		json.Unmarshal(b, &history)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return history, nil
}

// setupDB
func setupDB() (*bolt.DB, error) {
	db, err := bolt.Open("gronit.db", 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}
	return db, nil
}
