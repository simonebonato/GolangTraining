package main

import (
	"fmt"
	"net/http"
)


func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Welcome to Pippa's website!</h1>")
}

func main() {
	// if we put only the "/" it is going to handle all the requests that are made
	// this is to register the function so Go knows to use this when this type of request is received
	http.HandleFunc("/", handlerFunc)

	fmt.Println("Starting the server on :3000...")

	// this tells go to start a server and listen to port :3000
	// we can also give a specific domain 
	// if we put nil, it uses a default mux server
	http.ListenAndServe(":3000", nil)
}