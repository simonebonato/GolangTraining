package main

import (
	cyoa "adventure"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

func main() {
	// flags arguments
	filename := flag.String("file", "gopher.json", "The path to the JSON file with the CYOA story")
	port := flag.Int("port", 3000, "The port where the CYOA will run")
	flag.Parse()

	fmt.Printf("Using the story in %s.\n", *filename)

	// reading the file
	f, err := os.Open(*filename)
	if err != nil {panic(err)}

	story, err := cyoa.JSONStory(f)
	if err != nil {panic(err)}

	tpl := template.Must(template.New("").Parse("Hello"))
	story_handler := cyoa.NewHandler(story, cyoa.WithTemplate(tpl))

	fmt.Printf("Starting the server on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), story_handler))
}