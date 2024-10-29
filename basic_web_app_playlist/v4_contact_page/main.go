package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to SimoPippa's website!</h1>")
	fmt.Fprint(w, r.URL.Path)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hi, this it the contact page...</h1>")
	fmt.Fprint(w, "<p>To get in touch with me, send an email to: pippa@gmail.com\n")
	fmt.Fprint(w, r.URL.Path)
}

func main() {
	
	http.HandleFunc("/", handlerFunc)

	// we can register a new address and the HandlerFunc used to treat it
	// in this case it's another HTML but in the future we can do different things :)
	http.HandleFunc("/contact", contactHandler)

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", nil)
}