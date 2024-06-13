package main

import (
	"log"
	"net/http"
	"time"

	"oauth2-server/internal/config"

	"oauth2-server/internal/handler"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("error loading config: " + err.Error())
	}

	h := handler.New(*cfg)

	server := &http.Server{
		Addr:           ":8080",
		Handler:        h.Routes(),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   5 * time.Second,
	}

	log.Println("Starting server on :8080")
	log.Fatal(server.ListenAndServe())
}
