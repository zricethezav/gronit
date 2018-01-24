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

type Created struct {
	CreatedAt time.Time `json:"created_at"`
}

type Summary struct {
	StatusCount             int       `json:"status_count"`
	RunCount                int       `json:"run_count"`
	CompletionCount         int       `json:"completion_count"`
	AverageTimeToCompletion int64     `json:"average_time_to_completion"`
	CreatedAt               time.Time `json:"created_at"`
}

// initEntry create entry for new job
func initEntry(id string, db *bolt.DB, date time.Time) error {
	var err error
	createdAt := Created{CreatedAt: date}
	createdAtBytes, err := json.Marshal(createdAt)
	err = db.Update(func(tx *bolt.Tx) error {
		jobBucket, err := tx.CreateBucketIfNotExists([]byte(id))
		if err != nil {
			return fmt.Errorf("could not bucket for %s: %v", id, err)
		}
		err = jobBucket.Put([]byte("created"), createdAtBytes)
		if err != nil {
			return fmt.Errorf("failed to update status for %s: %v", id, err)
		}
		return nil
	})
	return err
}

// setStatus sets the state of the job
func setStatus(id string, statusStr string, date time.Time, db *bolt.DB) error {
	entry := Entry{Status: statusStr, Time: date}
	statusBytes, err := json.Marshal(entry)
	err = setDataBytes(id, "status", statusBytes, db)
	err = setHistory(id, &entry, db)
	return err
}

// setHistory updates the history of a job entry
func setHistory(id string, entry *Entry, db *bolt.DB) error {
	history, err := getHistory(id, db)
	history = append(history, *entry)
	historyBytes, err := json.Marshal(history)
	err = setDataBytes(id, "history", historyBytes, db)
	return err
}

// getStatus grabs the status of the job
func getStatus(id string, db *bolt.DB) (*Entry, error) {
	status := Entry{}
	statusBytes, err := getDataBytes(id, "status", db)
	json.Unmarshal(statusBytes, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// getHistory grabs the status of the job
func getHistory(id string, db *bolt.DB) ([]Entry, error) {
	history := []Entry{}
	historyBytes, err := getDataBytes(id, "history", db)
	json.Unmarshal(historyBytes, &history)
	if err != nil {
		return nil, err
	}
	return history, nil
}

// getSummary summarizes job entry stats
func getSummary(id string, db *bolt.DB) (*Summary, error) {
	var runTimeSum float64
	var previousEntry Entry
	var avg int64
	history, err := getHistory(id, db)
	createdBytes, err := getDataBytes(id, "created", db)
	created := Created{}
	json.Unmarshal(createdBytes, &created)
	if err != nil {
		return nil, err
	}
	completionCount := 0
	realCompletedJobCount := 0
	runCount := 0

	for _, entry := range history {
		if entry.Status == "running" {
			runCount += 1
		} else if entry.Status == "complete" {
			completionCount += 1
			if previousEntry.Status == "running" {
				runTimeSum += float64(entry.Time.Sub(previousEntry.Time))
				realCompletedJobCount += 1
			}
		}
		previousEntry = entry
	}

	if runCount != 0 {
		avg = int64(time.Duration(
			runTimeSum/float64(realCompletedJobCount)) / time.Millisecond)
	} else {
		avg = 0
	}

	summary := Summary{
		StatusCount:             len(history),
		RunCount:                runCount,
		CompletionCount:         completionCount,
		CreatedAt:               created.CreatedAt,
		AverageTimeToCompletion: avg,
	}
	return &summary, nil
}

// getDataBytes is a generic data aquisition helper
func getDataBytes(id string, bucket string, db *bolt.DB) ([]byte, error) {
	var data []byte
	err := db.View(func(tx *bolt.Tx) error {
		jobBucket := tx.Bucket([]byte(id))
		if jobBucket == nil {
			return fmt.Errorf("failed to find record of %s", id)
		}
		data = jobBucket.Get([]byte(bucket))
		return nil
	})
	return data, err
}

// setDataBytes generic set to bucket
func setDataBytes(id string, bucket string, data []byte, db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		jobBucket := tx.Bucket([]byte(id))
		if jobBucket == nil {
			return fmt.Errorf("failed to find record of %s", id)
		}
		err := jobBucket.Put([]byte(bucket),
			data)
		if err != nil {
			return fmt.Errorf("failed to update history for  %s: %v", id, err)
		}
		return nil
	})
	return err
}

// deletes bucket
func deleteBucket(id string, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(id))
	})
}

// setupDB
func setupDB() (*bolt.DB, error) {
	db, err := bolt.Open("gronit.db", 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}
	return db, nil
}
