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
	isbn := r.PathValue("isbn")
	book, err := BookRepo.GetByISBN(isbn)
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
	author := r.PathValue("author")
	books, err := BookRepo.GetByAuthor(author)
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

func GetInRange(w http.ResponseWriter, r *http.Request) {
	L.Info("GET /api/v1/books/range")
	Url, _ := url.Parse(r.URL.String())
	params, _ := url.ParseQuery(Url.RawQuery)
	year1, _ := strconv.Atoi(params["year1"][0])
	year2, _ := strconv.Atoi(params["year2"][0])

	books, err := BookRepo.GetInRange(year1, year2)
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
	_ = json.Unmarshal([]byte(string(body)), &bookData)
	for _, data := range bookData {
		existingBook, err := BookRepo.GetByISBN(data.ISBN)
		if err != nil {
			L.Error("Error GetByISBN: ", err)

			response := &Response{Status: "fail", Message: err.Error()}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		if existingBook == (repo.Book{}) {
			return
		}

		//
		res, err2 := BookRepo.Update(data)
		if err2 != nil {
			L.Error("Error Update: ", err2)
			response := &Response{Status: "fail", Message: err2.Error()}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		} else {
			response := &Response{Status: "success", Message: res}
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
			L.Error("Error GetByISBN: ", err)

			response := &Response{Status: "fail", Message: err.Error()}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		if existingBook == (repo.Book{}) {
			return
		}

		res, err2 := BookRepo.Delete(data)
		if err2 != nil {
			L.Error("Error Delete: ", err2)
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			response := &Response{Status: "fail", Message: err2.Error()}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		} else {
			response := &Response{Status: "success", Message: res}
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
		fmt.Print(data)
		_, err := BookRepo.GetByISBN(data.ISBN)
		if err == nil {
			L.Error("Error GetByISBN: ", err)

			response := &Response{Status: "fail", Message: ""}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		res, err2 := BookRepo.Insert(data)
		if err2 != nil {
			L.Error("Error Insert: ", err2)
			response := &Response{Status: "fail", Message: err2.Error()}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		} else {
			response := &Response{Status: "success", Message: res}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}
