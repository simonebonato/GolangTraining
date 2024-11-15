package boltDb

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/asdine/storm"
)

type Task struct {
	Key   int `storm:"id,increment"`
	Value string
}

type TaskSlice []Task

func (t TaskSlice) String() string {
	printString := "\nHere are your tasks for today:\n\n"
	for _, task := range t {
		printString += fmt.Sprintf("     --> %v  : %v\n", task.Key, task.Value)
	}

	return printString
}

// bolt uses byte slices as keys to access data
// so we need to convert the integer values of the todo to byte slices
// also the data has to be decoded into byte slices (eg. canno throw directly a struct to it)
// in the end I use storm, which is built on top of bolt and makes life 1000000 times easier :D
func CreateDb(folderName string) (refToDb *storm.DB) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	newPath := filepath.Join(cwd, folderName, "todos.db")

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := storm.Open(newPath)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func AddTask(refToDb *storm.DB, newTask *Task) {
	err := refToDb.Save(newTask)
	if err != nil {
		panic(err)
	}
}

func ResetTaskKeys(refToDb *storm.DB) error {
	var allTasks TaskSlice
	err := refToDb.All(&allTasks)
	if err != nil {
		return err
	}

	err = refToDb.Select().Delete(new(Task))
	if err != nil {
		fmt.Println("Error deleting entries for the sorting!")
	}

	// Reassign Keys and update tasks
	for i, task := range allTasks {
		task.Key = i + 1
		err = refToDb.Save(&task)
		if err != nil {
			return err
		}
	}

	return nil

}
