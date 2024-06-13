package handler

import (
	"crypto/rsa"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"
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

type Handler struct {
	mux *http.ServeMux
}

func New() *Handler {
	return &Handler{
		mux: http.NewServeMux(),
	}
}

func (h *Handler) Routes() http.Handler {
	h.mux.HandleFunc("POST /token", h.generateToken)
	protectedEndpointHandler := http.HandlerFunc(h.protectedEndpoint)
	h.mux.Handle("GET /protected", h.validateTokenMiddleware(protectedEndpointHandler))

	return h.mux
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (h *Handler) generateToken(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) protectedEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a protected endpoint"))
}
