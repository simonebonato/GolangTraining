package boltDb

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
)

var taskBucket = []byte("tasks")

type Task struct {
	Key int
	Value  string
}

// bolt uses byte slices as keys to access data
// so we need to convert the integer values of the todo to byte slices
// also the data has to be decoded into byte slices (eg. canno throw directly a struct to it)
func CreateDb(folderName string) (refToDb *bolt.DB){
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	newPath := filepath.Join(cwd, folderName, "todos.db")
	fmt.Println(newPath)

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(newPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// add the bucket to the database, like adding a model
	err = db.Update(func (tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err
		})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
