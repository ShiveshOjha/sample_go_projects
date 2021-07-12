package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux" // for routing
)

type Book struct {
	ID     int    `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   string `json:year`
}

var books []Book

func main() {

	router := mux.NewRouter()

	books = append(books,
		Book{ID: 1, Title: "Golang pointers", Author: "Mr. Golang", Year: "2012"},
		Book{ID: 2, Title: "Golang routines", Author: "Mr. Goroutine", Year: "2011"},
	)

	router.HandleFunc("/books", getBooks).Methods("GET") // (endpoint, handlerFunction)
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8585", router))

}

func getBooks(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(books)

}

func getBook(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	k, _ := strconv.Atoi(params["id"])

	for _, book := range books {
		if book.ID == k {
			json.NewEncoder(w).Encode(&book)
		}
	}

}

func addBook(w http.ResponseWriter, r *http.Request) { //send json object to server and decode request body and map the values inside that body to the book

	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	books = append(books, book)

	json.NewDecoder(r.Body).Decode(books) // decoding request body and maping it to book values

}

func updateBook(w http.ResponseWriter, r *http.Request) { // put

	var book Book

	json.NewDecoder(r.Body).Decode(&book) // map attributes to Book

	for i, item := range books {
		if item.ID == book.ID {
			books[i] = book
		}
	}

	json.NewEncoder(w).Encode(books)

}

func removeBook(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	for i, item := range books {
		if item.ID == id {
			books = append(books[:i], books[i+1:]...) // make a new slice which doesn't have inputted id at index i
		}
	}

	json.NewEncoder(w).Encode(books)

}
