package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json: "id"`
	Isbn     string    `json: "isbn"`
	Title    string    `jsin: "title"`
	Director *Director `json: "director"`
}

type Director struct {
	Firstname string `json: "firstname"`
	Lastname  string `json: "lastname"`
}

var movies []Movie //list vars are of type Movie (The struct)

// Utility
func RemoveIndex(s []Movie, index int) []Movie {
	return append(s[:index], s[index+1:]...)
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies) // returning movies in json format
	return
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //setting content type to json (easier to work with)
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = RemoveIndex(movies, index)
			break
		}
	}
	json.NewEncoder(w).Encode(movies) //return remaining movies
	return
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item) // return needed movies
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie) // read and decode movie from front end
	movie.ID = strconv.Itoa(rand.Intn(100000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
	return
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = RemoveIndex(movies, index)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie) // return movie
			return
		}
	}
}

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie 1", Director: &Director{Firstname: "Director", Lastname: "One"}}) //& is to get the address (kinda the opposite of pointers)
	movies = append(movies, Movie{ID: "2", Isbn: "126543", Title: "Movie 2", Director: &Director{Firstname: "Director", Lastname: "Two"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting Server at Port 8000 \n")
	log.Fatal(http.ListenAndServe(":8000", r))

}
