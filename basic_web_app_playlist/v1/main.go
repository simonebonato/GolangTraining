package main

import (
	"fmt"
	"net/http"
)

// the request is, well, the request that is made by the server, information that is sent to us
// the response writer instead is a way the server can reply back to the server

// ResponseWriter is an INTERFACE, because we can reply with different types of Writers
// while the Reqest is a struct, since normally it has to follow a specific format

// so functions that match this definition can be used to handle web requests in GO
// there is a TYPE defining it called HandlerFunc, and we want to match it
// we can call it however but the definition must be the same
func handlerFunc(w http.ResponseWriter, r *http.Request) {
	// fprint let's us choose where we can write to
	fmt.Fprint(w, "<h1>Welcome to Pippa's website!</h1>")
}

func main() {
	http.HandleFunc("/", handlerFunc)
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", nil)
}