package main

import (
	"encoding/json"
	"fmt"
	"github.com/ladnaaaaaa/calc_service/internal/calculator"
	"log"
	"net/http"
	"strings"
	"unicode"
)

type CalculateRequest struct {
	Expression string `json:"expression"`
}

type CalculateResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var request CalculateRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		respondWithError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if !validateExpression(w, request) {
		respondWithError(w, "Expression is not valid", http.StatusUnprocessableEntity)
		return
	}
	result, err := calculator.Calc(request.Expression)
	if err != nil {
		respondWithError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondWithJson(w, CalculateResponse{Result: fmt.Sprint(result)}, http.StatusOK)
}

func validateExpression(w http.ResponseWriter, request CalculateRequest) bool {
	expression := strings.ReplaceAll(request.Expression, " ", "")
	for _, char := range expression {
		if !unicode.IsDigit(char) && !strings.ContainsRune("+-*/()", char) {
			return false
		}
	}
	return true
}

func respondWithError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(CalculateResponse{Error: message})
}

func respondWithJson(w http.ResponseWriter, payload CalculateResponse, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Could not start server: %v\n ", err)
	}
}
