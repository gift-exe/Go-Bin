package main

import (
	"fmt"
	"log"
	"net/http"
)

func formHandler(w http.ResponseWriter, r *http.Request) { //'*' is a pointer
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "POST request successful")
	name := r.FormValue("name")
	address := r.FormValue("address")
	fmt.Fprintf(w, "Name: %v\n", name)
	fmt.Fprintf(w, "Address: %v\n", address)

}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Hello !")
}

func main() {
	fileServer := http.FileServer(http.Dir("./static")) //gets the index page from static dir (generally the first page)
	http.Handle("/", fileServer)                        //handles the index page
	http.HandleFunc("/form", formHandler)               //handles the form page
	http.HandleFunc("/hello", helloHandler)             //handles the hello page

	//inference:= handling simply involves connecting functions to static pages

	fmt.Printf("Server Starting at port 8080 \n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
