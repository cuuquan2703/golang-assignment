package routers

import (
	"encoding/json"
	"fmt"
	_ "fmt"
	"io"
	"net/http"
	"net/url"
	"server/logger"
	repo "server/repositories"
	"strconv"
)

type Response struct {
	Status  string `json:"status"`
	Message any    `json:"message"`
}

var BookRepo, _ = repo.NewBookRepository()
var L = logger.CreateLog()

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	L.Info("GET /api/v1/books")
	books, err := BookRepo.GetAllBooks()
	if err != nil {
		L.Error("Error: ", err)
		response := &Response{Status: "fail", Message: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	response := &Response{Status: "success", Message: books}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetByISBN(w http.ResponseWriter, r *http.Request) {
	L.Info("GET /api/v1/books/i/")
	isbn := r.PathValue("isbn")
	book, err := BookRepo.GetByISBN(isbn)
	if err != nil {
		L.Error("Error: ", err)
		response := &Response{Status: "fail", Message: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	response := &Response{Status: "success", Message: book}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetByAuthor(w http.ResponseWriter, r *http.Request) {
	L.Info("GET /api/v1/books/a/")
	author := r.PathValue("author")
	books, err := BookRepo.GetByAuthor(author)
	if err != nil {
		L.Error("Error: ", err)
		response := &Response{Status: "fail", Message: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	response := &Response{Status: "success", Message: books}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetInRange(w http.ResponseWriter, r *http.Request) {
	L.Info("GET /api/v1/books/range")
	Url, _ := url.Parse(r.URL.String())
	params, _ := url.ParseQuery(Url.RawQuery)
	year1, _ := strconv.Atoi(params["year1"][0])
	year2, _ := strconv.Atoi(params["year2"][0])

	books, err := BookRepo.GetInRange(year1, year2)
	if err != nil {
		L.Error("Error: ", err)
		response := &Response{Status: "fail", Message: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	response := &Response{Status: "success", Message: books}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Update(w http.ResponseWriter, r *http.Request) {
	L.Info("POST /api/v1/books/update")
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	var bookData []repo.Book
	_ = json.Unmarshal([]byte(string(body)), &bookData)
	for _, data := range bookData {
		existingBook, err := BookRepo.GetByISBN(data.ISBN)
		if err != nil {
			L.Error("Error: ", err)

			response := &Response{Status: "fail", Message: err.Error()}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		if existingBook == (repo.Book{}) {
			return
		}

		_, err2 := BookRepo.DB.Exec("UPDATE Book SET name = $1, publish_year = $2, author = $3 WHERE isbn = $4", data.Name, data.PublishYear, data.Author, data.ISBN)
		if err2 != nil {
			L.Error("Error: ", err2)
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			response := &Response{Status: "fail", Message: err2.Error()}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	L.Info("DELETE /api/v1/books/delete")
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	var bookData []repo.Book
	_ = json.Unmarshal([]byte(string(body)), &bookData)
	for _, data := range bookData {
		existingBook, err := BookRepo.GetByISBN(data.ISBN)
		if err != nil {
			L.Error("Error: ", err)

			response := &Response{Status: "fail", Message: err.Error()}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		if existingBook == (repo.Book{}) {
			return
		}

		_, err2 := BookRepo.DB.Exec("DELETE FROM Book WHERE isbn = $1", data.ISBN)
		if err2 != nil {
			L.Error("Error: ", err2)
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			response := &Response{Status: "fail", Message: err2.Error()}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}

func Insert(w http.ResponseWriter, r *http.Request) {
	L.Info("POST /api/v1/books/add")
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	var bookData []repo.Book
	_ = json.Unmarshal([]byte(string(body)), &bookData)
	fmt.Println("Request Body:", bookData)
	for _, data := range bookData {
		_, err := BookRepo.GetByISBN(data.ISBN)
		if err == nil {
			L.Error("Error: ", err)

			response := &Response{Status: "fail", Message: ""}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		_, err2 := BookRepo.DB.Exec(" INSERT INTO Book (isbn, name, publish_year, author) VALUES ($1, $2, $3, $4);", data.ISBN, data.Name, data.PublishYear, data.Author)
		if err2 != nil {
			L.Error("Error: ", err2)
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			response := &Response{Status: "fail", Message: err2.Error()}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}
