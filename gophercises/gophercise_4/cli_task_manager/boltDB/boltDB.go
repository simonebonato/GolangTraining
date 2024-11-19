package boltDb

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

type KeyHolder interface {
	GetKey() int
}

type Task struct {
	Key   int `storm:"id,increment"`
	Value string
}

func (t *Task) GetKey() int {
	return t.Key
}

type DoneTask struct {
	Id       int `storm:"id,increment"`
	DoneTask Task
	DoneAt   time.Time `storm:"index"`
}

func (dt *DoneTask) GetKey() int {
	return dt.DoneTask.Key
}

func (dt *DoneTask) Validate() error {
	if dt.DoneTask.Value == "" {
		return errors.New("DoneTask.Value must not be empty")
	}
	return nil
}

type TaskSlice []Task

func (t TaskSlice) String() string {
	printString := "\nHere are your tasks for today:\n\n"
	for _, task := range t {
		printString += fmt.Sprintf("     --> %v  : %v\n", task.Key, task.Value)
	}

	return printString
}

type DoneTaskSlice []DoneTask

func (dt DoneTaskSlice) String() string {
	printString := "\nHere are all the tasks you did so far:\n\n"
	for _, task := range dt {
		printString += fmt.Sprintf("     --> %v  : %v\n", task.DoneTask.Value, task.DoneAt)
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

func AddTask(refToDb *storm.DB, newTask KeyHolder) {
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

	if len(allTasks) == 0 {
		return nil
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

func PrintToDos(refToDb *storm.DB) {
	var tasks TaskSlice
	err := refToDb.All(&tasks)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tasks)
}

func PrintDoneToDos(refToDb *storm.DB) {
	var doneTasks DoneTaskSlice
	err := refToDb.Select(q.True(), q.Gt("DoneAt", time.Now().Add(-24*time.Hour))).OrderBy("DoneAt").Reverse().Find(&doneTasks)
	if len(doneTasks) == 0 {
		fmt.Println("There are no tasks here.")
		return
	} else if err != nil {
		log.Fatal(err)
	}
	fmt.Println(doneTasks)
}
