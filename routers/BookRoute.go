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
var L = logger.CreateLog()
var bookService = service.BookService{
	Repo: BookRepo,
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	L.Info("GET /api/v1/books")
	books, err := bookService.GetAllBooks()
	if err != nil {
		L.Error("Error GetAllBooks: ", err)
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
	// isbn := r.PathValue("isbn")
	Url, _ := url.Parse(r.URL.String())
	params, _ := url.ParseQuery(Url.RawQuery)
	book, err := bookService.GetByISBN(params["isbn"][0])
	if err != nil {
		L.Error("Error GetByISBN: ", err)
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
	// author := r.PathValue("author")
	Url, _ := url.Parse(r.URL.String())
	params, _ := url.ParseQuery(Url.RawQuery)
	books, err := bookService.GetByAuthor(params["author"][0])
	if err != nil {
		L.Error("Error GetByAuthor: ", err)
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

func GetInRange(w http.ResponseWriter, r *http.Request) {
	L.Info("GET /api/v1/books/range")
	Url, _ := url.Parse(r.URL.String())
	params, _ := url.ParseQuery(Url.RawQuery)
	fmt.Println(params)
	from, _ := strconv.Atoi(params["from"][0])
	to, _ := strconv.Atoi(params["to"][0])

	books, err := bookService.GetInRange(from, to)
	if err != nil {
		L.Error("Error in range: ", err)
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
	er := json.Unmarshal([]byte(string(body)), &bookData)
	if er != nil {
		L.Error("json.Unmarshal", er)
	}
	res, err := bookService.Update(bookData)
	if err != nil {
		response := &Response{Status: "fail", Message: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		response := &Response{Status: "success", Message: res}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	L.Info("DELETE /api/v1/books/delete")
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	var bookData []repo.Book
	er := json.Unmarshal([]byte(string(body)), &bookData)
	if er != nil {
		L.Error("json.Unmarshal", er)
	}
	res, err := bookService.Delete(bookData)
	if err != nil {
		response := &Response{Status: "fail", Message: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		response := &Response{Status: "success", Message: res}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func Insert(w http.ResponseWriter, r *http.Request) {
	L.Info("POST /api/v1/books/add")
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	var bookData []repo.Book
	er := json.Unmarshal([]byte(string(body)), &bookData)
	if er != nil {
		L.Error("json.Unmarshal", er)
	}
	res, err := bookService.Insert(bookData)
	if err != nil {
		response := &Response{Status: "fail", Message: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		response := &Response{Status: "success", Message: res}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
