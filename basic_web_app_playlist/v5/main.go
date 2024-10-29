package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to SimoPippa's website!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hi, this is the contact page...</h1>")
	fmt.Fprint(w, "<p>To get in touch with me, send an email to: pippa@gmail.com\n")
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status) // this is to send back the status code to the browser!
	switch status {
	case http.StatusNotFound:
		http.Error(w, "Page not found, sorry :(", status)
	default:
		fmt.Fprint(w, "<h1>Page not found dude</h1>")
	}
}

// why do we need a struct of this sort instead of defining the functions like we did in the
// previous lessons? we can have access to different fields that have access to the struct
// and we can also mix other structs together
type Router struct {}

func (router Router) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch path {
	case "/":
		handlerFunc(w, r)
	case "/contact":
		contactHandler(w, r)
	default:
		errorHandler(w, r, http.StatusNotFound)
	}
}

func main() {
	var router Router
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", router)
}