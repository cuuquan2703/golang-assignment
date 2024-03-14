package routers

import (
	"encoding/json"
	_ "fmt"
	"net/http"
	"net/url"
	repo "server/repositories"
	"strconv"
)

// type BookRoute struct {
// 	repo *repo.BookRepository
// }

// const BookRepo = *repo.BookRepository

var BookRepo, _ = repo.NewBookRepository()

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := BookRepo.GetAllBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func GetByISBN(w http.ResponseWriter, r *http.Request) {
	isbn := r.PathValue("isbn")
	book, err := BookRepo.GetByISBN(isbn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func GetByAuthor(w http.ResponseWriter, r *http.Request) {
	author := r.PathValue("author")
	books, err := BookRepo.GetByAuthor(author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func GetInRange(w http.ResponseWriter, r *http.Request) {
	Url, _ := url.Parse(r.URL.String())
	params, _ := url.ParseQuery(Url.RawQuery)
	year1, _ := strconv.Atoi(params["year1"][0])
	year2, _ := strconv.Atoi(params["year2"][0])
	books, err := BookRepo.GetInRange(year1, year2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}
