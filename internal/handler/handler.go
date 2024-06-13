package handler

import (
	"crypto/rsa"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"oauth2-server/internal/config"

	"github.com/dgrijalva/jwt-go"
)

type Handler struct {
	clientID      string
	clientSecret  string
	rsaPublicKey  *rsa.PublicKey
	rsaPrivateKey *rsa.PrivateKey
}

type ProtectedResponse struct {
	Message string `json:"message"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func New(cfg config.Config, rsaPublicKey *rsa.PublicKey, rsaPrivateKey *rsa.PrivateKey) *Handler {
	return &Handler{
		clientID:      cfg.ClientID,
		clientSecret:  cfg.ClientSecret,
		rsaPublicKey:  rsaPublicKey,
		rsaPrivateKey: rsaPrivateKey,
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /token", h.generateToken)
	protectedEndpointHandler := http.HandlerFunc(h.protectedEndpoint)
	mux.Handle("GET /protected", h.validateTokenMiddleware(protectedEndpointHandler))

	return mux
}

func (h *Handler) generateToken(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handleError(w, http.StatusBadRequest, "invalid request")

		return
	}
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok || clientID != h.clientID || clientSecret != h.clientSecret {
		handleError(w, http.StatusUnauthorized, "invalid credentials")

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "some_server_id",
		"sub": clientID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(h.rsaPrivateKey)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err.Error())

		return
	}

	response := TokenResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) protectedEndpoint(w http.ResponseWriter, _ *http.Request) {
	msg := ProtectedResponse{
		Message: "This is a protected endpoint",
	}
	resp, err := json.Marshal(&msg)
	if err != nil {
		log.Printf("protected Marshal error: %v\n", err)

		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
