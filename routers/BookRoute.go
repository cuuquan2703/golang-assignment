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
	// isbn := r.PathValue("isbn")
	Url, _ := url.Parse(r.URL.String())
	params, _ := url.ParseQuery(Url.RawQuery)
	book, err := BookRepo.GetByISBN(params["isbn"][0])
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
	books, err := BookRepo.GetByAuthor(params["author"][0])
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

	books, err := BookRepo.GetInRange(from, to)
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
		_, err := BookRepo.GetByISBN(data.ISBN)
		if err != nil {
			L.Error("Error GetByISBN: ", err)

			response := &Response{Status: "fail", Message: err.Error()}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
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
		}
	}
}
