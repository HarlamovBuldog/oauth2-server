package main

import (
	"context"
	"crypto/rsa"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dgrijalva/jwt-go"

	"oauth2-server/internal/config"
	"oauth2-server/internal/handler"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func init() {
	privateKeyData, err := os.ReadFile("private.key")
	if err != nil {
		panic(err)
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		panic(err)
	}

	publicKeyData, err := os.ReadFile("public.key")
	if err != nil {
		panic(err)
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		panic(err)
	}
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("error loading config: " + err.Error())
	}

	h := handler.New(*cfg, publicKey, privateKey)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        h.Routes(),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   5 * time.Second,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Println("Gracefully shitting down the http server")

				return
			}

			log.Printf("failed to run the http server: %v\n", err.Error())
		}
	}()

	log.Println("Starting server on :8080")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shut down the http server: %v\n", err.Error())
	}
}
