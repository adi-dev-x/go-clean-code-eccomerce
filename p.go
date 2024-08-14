package main

import (
	"fmt"
	"log"
	"net/http"
	// "os"
)

func main() {
	// Define the file server to serve files from the current directory
	fileServer := http.FileServer(http.Dir("./"))

	// Serve files from the root URL
	http.Handle("/", fileServer)

	// Define the port on which the server will listen
	port := "8081"
	fmt.Printf("Starting server at http://localhost:%s/\n", port)

	// Start the server
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
