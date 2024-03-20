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
	"server/service"
	"strconv"
)

type Response struct {
	Status  string `json:"status"`
	Message any    `json:"message"`
}

var BookRepo, _ = repo.NewBookRepository()
var BookService = service.BookService{
	Repo: BookRepo,
}
var L = logger.CreateLog()

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	L.Info("GET /api/v1/books")
	books, err := BookService.GetAllBooks()
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
	L.Info("GET /api/v1/books?{isbn}")
	Url, _ := url.Parse(r.URL.String())
	params, _ := url.ParseQuery(Url.RawQuery)
	book, err := BookService.GetByISBN(params["isbn"][0])
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
	L.Info("GET /api/v1/books?{author}")
	Url, _ := url.Parse(r.URL.String())
	params, _ := url.ParseQuery(Url.RawQuery)
	books, err := BookService.GetByAuthor(params["author"][0])
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
	from, _ := strconv.Atoi(params["from"][0])
	to, _ := strconv.Atoi(params["to"][0])

	books, err := BookService.GetInRange(from, to)
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

func Get(w http.ResponseWriter, r *http.Request) {
	Url, _ := url.Parse(r.URL.String())
	params, _ := url.ParseQuery(Url.RawQuery)
	if len(params) == 0 {
		GetAllBooks(w, r)
	}
	if _, ok := params["isbn"]; ok {
		GetByISBN(w, r)
	}
	if _, ok := params["author"]; ok {
		GetByAuthor(w, r)
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	L.Info("POST /api/v1/books/update")
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	var bookData []repo.Book
	_ = json.Unmarshal([]byte(string(body)), &bookData)
	for _, data := range bookData {
		existingBook, err := BookService.GetByISBN(data.ISBN)
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

		_, err2 := BookService.Update(data.ISBN, data.Name, data.Author, data.PublishYear)
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
		existingBook, err := BookService.GetByISBN(data.ISBN)
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

		_, err2 := BookService.Delete(data.ISBN)
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
		_, err := BookService.GetByISBN(data.ISBN)
		if err == nil {
			L.Error("Error: ", err)

			response := &Response{Status: "fail", Message: ""}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		_, err2 := BookService.Insert(data.ISBN, data.Name, data.Author, data.PublishYear)
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
