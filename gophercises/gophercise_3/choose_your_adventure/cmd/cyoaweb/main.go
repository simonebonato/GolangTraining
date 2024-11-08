package main

import (
	cyoa "adventure"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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

	

	tpl := template.Must(template.New("").Parse(cyoa.StoryTmpl))
	// this way we can use the 2 custom options for the story handling
	story_handler := cyoa.NewHandler(story, cyoa.WithPathFn(newPathFn), cyoa.WithTemplate(tpl))

	// we can make this ServeMux to handle all the requests that start with /story/ to go to the story handler
	// while all the rest will get the 404 error page
	mux := http.NewServeMux()
	mux.Handle("/story/", story_handler)

	fmt.Printf("Starting the server on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}


func newPathFn(r *http.Request) string {
	url := strings.TrimSpace(r.URL.Path)

	var render_chapter string
	if url == "/story" || url == "/story/" {
		render_chapter = "intro"
	} else {
		render_chapter = url[len("/story/"):]
	}
	return render_chapter
}