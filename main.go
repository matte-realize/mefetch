package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"mefetch/handlers"
)

func main() {
	godotenv.Load()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.Health)

	port := os.Getenv("PORT")

	fmt.Printf("Server running on http://localhost:%s\n", port)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}