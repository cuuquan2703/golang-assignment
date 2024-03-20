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
	http.HandleFunc("GET /api/v1/books", Route.Get)
	http.HandleFunc("GET /api/v1/books/range", Route.GetInRange)
	http.HandleFunc("POST /api/v1/books/update", Route.Update)
	http.HandleFunc("DELETE /api/v1/books/delete", Route.Delete)
	http.HandleFunc("POST /api/v1/books/add", Route.Insert)

	http.ListenAndServe(":8081", nil)
}
