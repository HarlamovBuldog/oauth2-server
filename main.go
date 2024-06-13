package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
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

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func generateToken(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok || clientID != "client-id" || clientSecret != "client-secret" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "some_server_id",
		"sub": clientID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := TokenResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\n")
	if err := encoder.Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		panic("Client ID or Client Secret not set in .env file")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /token", generateToken)
	protectedEndpointHandler := http.HandlerFunc(protectedEndpoint)
	mux.Handle("GET /protected", validateTokenMiddleware(protectedEndpointHandler))

	server := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   5 * time.Second,
	}

	log.Println("Starting server on :8080")
	log.Fatal(server.ListenAndServe())
}

func validateTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return publicKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func protectedEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a protected endpoint"))
}
