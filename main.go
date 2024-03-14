package main

import (
	"net/http"

	// r "server/repositories"
	Route "server/routers"
)

func main() {

	// defer db.Close()
	// if err := insertMockData(db); err != nil {
	// 	log.Fatal(err)
	// }
	http.HandleFunc("GET /api/v1/books", Route.GetAllBooks)
	http.HandleFunc("GET /api/v1/books/{isbn}", Route.GetByISBN)
	http.HandleFunc("GET /api/v1/books/{author}", Route.GetByAuthor)
	http.HandleFunc("GET /api/v1/books/range", Route.GetInRange)

	http.ListenAndServe(":8081", nil)
}
