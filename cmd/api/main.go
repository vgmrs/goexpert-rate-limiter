package main

import (
	"log"
	"net/http"

	"github.com/vgmrs/goexpert-rate-limiter/internal/config"
)

func main() {
	rateLimiterMiddleware, err := config.SetupRateLimiter()
	if err != nil {
		log.Fatalf("Rate limiter error: %v", err)
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Rate Limiter"))
	})

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, rateLimiterMiddleware(handler)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
