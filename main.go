package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"
)

type Book struct {
	Author        string `json:"author"`
	Title         string `json:"title"`
	YearPublished uint   `json:"yearPublished"`
	Id            uint   `json:"id"`
}

var books []Book  // Simple in-memory storage, obviously not persistent

type BookResponse struct {
	Books []Book `json:"books"`
}

// curl --request GET http://localhost:4000/health
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Don't panic.")
}

/*
curl --data '{"author": "Douglas Adams", "title": "The Hitchhikers Guide to the Galaxy", "yearPublished": 1979}' --header "Content-Type: application/json" --request POST http://localhost:4000/api/books

curl --data '{"title": "Moby Dick", "author": "Herman Melville", "yearPublished": 1851}' --header  "Content-Type: application/json" --request POST http://localhost:4000/api/books

curl --data '{"author": "Philip K. Dick", "title": "Do Androids Dream of Electric Sheep?", "yearPublished": 1968}' --header  "Content-Type: application/json" --request POST http://localhost:4000/api/books

*/
func AddBookHandler(w http.ResponseWriter, r *http.Request) {
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	var b Book
	var unmarshallErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&b)
	if err != nil {
		if errors.As(err, &unmarshallErr) {
			errorResponse(w, fmt.Sprintf("bad request:  incorrect type provided for field %s", unmarshallErr.Field), http.StatusBadRequest)
		} else {
			errorResponse(w, fmt.Sprintf("bad request: %s", err.Error()), http.StatusBadGateway)
		}
	}

	// Update the book with an ID:
	b.Id = (uint)(len(books) + 1)

	// Add the book to in-memory storage:
	books = append(books, b)

	jsonTxt, err := json.Marshal(b)
	if err != nil {
		errorResponse(w, fmt.Sprintf("json marshall error: %s", err.Error()), http.StatusInternalServerError)
	}

	// Write response:
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonTxt)
}

// curl --request DELETE http://localhost:4000/api/books  --verbose
func DeleteBooksHandler(w http.ResponseWriter, r *http.Request) {
	books = books[:0]
	w.WriteHeader(http.StatusNoContent)
}

func sortByTitle(a int, b int) bool {
	return books[a].Title < books[b].Title
}

// curl --request GET http://localhost:4000/api/books
func GetBookHandler(w http.ResponseWriter, r *http.Request) {
	// Sort the books by title:
	sort.Slice(books[:], sortByTitle)

	// Create the response data structure:
	bookResponse := BookResponse{books}

	jsonTxt, err := json.Marshal(bookResponse)
	if err != nil {
		errorResponse(w, fmt.Sprintf("json marshall error: %s", err.Error()), http.StatusInternalServerError)
	}

	// Write response:
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonTxt)
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", HealthHandler).Methods("GET")
	r.HandleFunc("/api/books", AddBookHandler).Methods("POST")
	r.HandleFunc("/api/books", GetBookHandler).Methods("GET")
	r.HandleFunc("/api/books", DeleteBooksHandler).Methods("DELETE")
	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:4000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
