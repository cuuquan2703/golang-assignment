package main

import (
    "net/http"
	routes"server/routers"
)

func main() {
    http.HandleFunc("/book",routes.Handler)

    // Start the server
    http.ListenAndServe(":8080", nil)
}