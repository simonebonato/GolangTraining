package main

import (
	"flag"
	"fmt"
	"net/http"
	"urlshort"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	// load a yaml as a flag, parse it and use it
	yaml_path := flag.String("yaml_path", "", "Insert the path to a yaml file")
	json_path := flag.String("json_path", "", "Insert the path to a json file")

	flag.Parse()

	yamlFILEHandler, err := urlshort.YAMLFileHandler(*yaml_path, yamlHandler)
	if err != nil {
		panic(err)
	}

	jsonFILEHandler, err := urlshort.JSONFIleHandler(*json_path, yamlFILEHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonFILEHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
