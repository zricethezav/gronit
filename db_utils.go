package main

import (
	_ "flag"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

func set(token string, data []byte) error {
	// filt out make generic?
	var err error
	err = db.Update(func(tx *bolt.Tx) error {
		err = tx.Bucket([]byte("token")).Put([]byte(token), data)
		if err != nil {
			return fmt.Errorf("could not insert entry: %v", err)
		}

		return nil
	})
	fmt.Println("Added Entry")
	return nil
}

func listBucket(db *bolt.DB) error {
	// TODO make getPath function that accepts a slice of variable length, then traverses it until it gets the data it needs
	var err error
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("token"))
		b.ForEach(func(k, v []byte) error {
			fmt.Println(string(k), string(v))
			return nil
		})
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func setupDB() (*bolt.DB, error) {
	db, err := bolt.Open("gronit.db", 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("token"))
		if err != nil {
			return fmt.Errorf("could not token bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not set up buckets, %v", err)
	}
	fmt.Println("DB Setup Done")
	return db, nil
}
