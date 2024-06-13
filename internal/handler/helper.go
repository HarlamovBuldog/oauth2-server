package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type Error struct {
	Error string `json:"error"`
}

func errorBody(message string) ([]byte, error) {
	resp := Error{
		Error: message,
	}

	return json.Marshal(resp)
}

func handleError(w http.ResponseWriter, status int, errMsg string) {
	resp, err := errorBody(errMsg)
	if err != nil {
		log.Printf("handleError NewErrorBody error: %v\n", err.Error())

		http.Error(w, errMsg, status)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)
}
