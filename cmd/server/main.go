package main

import (
	"log"
	"net/http"
	"os"

	"github.com/torbenish/clima-cep/internal/cep"
	"github.com/torbenish/clima-cep/internal/handler"
	"github.com/torbenish/clima-cep/internal/weather"
)

func main() {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	h := &handler.Handler{
		CEP:     cep.NewClient(),
		Weather: weather.NewClient(apiKey),
	}

	mux := http.NewServeMux()
	mux.Handle("GET /{cep}", h)

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}