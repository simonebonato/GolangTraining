package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	// we can use headers to influence the behaviour of what we send to the server
	// this will change how the content is interpreted by the browser
	// w.Header().Set("Content-Type", "text/plain")

	// how to figure out how to change the header set?
	// need to find the browser nomenclature online, not specific to Go

	// how to figure out how to set stuff up? The "Content-type" for isntance?
	// go to the docs of the Response writer and see the available methods
	// so we can see all we can do with the ResponseWriter 
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to SimoPippa's website!</h1>")
}

func main() {
	
	http.HandleFunc("/", handlerFunc)

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", nil)
}