package main

import (
	"log"

	"github.com/abhinav1912/commit-log/internal/server"
)

func main() {
	server := server.NewHTTPServer(":8080")
	log.Fatal(server.ListenAndServe())
}
