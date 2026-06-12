package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"mefetch/handlers"
)

func main() {
	godotenv.Load()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/generate", handlers.Generate)
	mux.HandleFunc("/card.svg", handlers.CardLive)
	mux.HandleFunc("/ascii", handlers.CardLive)
	mux.HandleFunc("/card/generate", handlers.CardGenerate)
	mux.Handle("/", http.FileServer(http.Dir("static")))

	port := os.Getenv("PORT")

	go func() {
		time.Sleep(500 * time.Millisecond)
		openBrowser("http://localhost:" + port)
	}()

	fmt.Printf("Server running on http://localhost:%s\n", port)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		fmt.Println("could not detect OS - open your browser manually at", url)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("could not open browser:", err)
	}
}