// link to exercise: https://courses.calhoun.io/lessons/les_goph_01
package main

import (
	"fmt"
	"os"
	"log"
	"encoding/csv"
	"github.com/akamensky/argparse"
	"time"
)

type QuizScore struct {
	correct_answers int8
	total_questions int8
}

func (quiz *QuizScore) answered_correctly () {
	quiz.correct_answers += 1
}

func (quiz *QuizScore) add_question_counter () {
	quiz.total_questions += 1
}

func (quiz QuizScore) String() string {
	return fmt.Sprintf("\nYou answered correctly %d question out of %d", quiz.correct_answers, quiz.total_questions)
}


func main() {
	parser := argparse.NewParser("quiz_parser", "This is a quiz game :)")

	filename := parser.String("f", "filename", &argparse.Options{Required: false, Help: "Name or path to the .csv file", Default: "problems.csv"})
	timelimit := parser.Int("l", "limit", &argparse.Options{Required: false, Help: "Time limit for the quiz."})

	err := parser.Parse(os.Args)
    if err != nil {
        fmt.Print(parser.Usage(err))
        return
    }

	file, err := os.Open(*filename)

	if err != nil {
		log.Printf("The file %v could not be loaded\n", *filename)
		return
		} else {log.Printf("The file %v has been loaded succesfully!\n", *filename)}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {log.Println("\nError reading the .csv file...")}

	var TheQuiz QuizScore 
	
	// start the timer function
	timer := time.NewTimer(time.Duration(*timelimit) * time.Second)

	for _, qa := range records {
		question, correct_answer := qa[0], qa[1]

		fmt.Println(question, "=")
		TheQuiz.add_question_counter()
		
		answerChan := make(chan string)
		go func() {
			var input string
			fmt.Scan(&input)
			answerChan <- input
		}()

		select {
		case <-timer.C:
			fmt.Println("\nTimer expired!")
			fmt.Println(TheQuiz)
			return
		case input := <- answerChan:
			if input == correct_answer {TheQuiz.answered_correctly()}
	}
}

	// make a stringer method for the quiz!
	fmt.Println(TheQuiz)

}