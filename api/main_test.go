package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateTransactionHandler_InvalidMethod(t *testing.T) {
	db = &mockDB{} // Use the mock so db is not nil

	req, err := http.NewRequest(http.MethodGet, "/api/transaction/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createTransactionHandler)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %v, want %v", rr.Code, http.StatusMethodNotAllowed)
	}
}

func TestCreateTransactionHandler_InvalidJSON(t *testing.T) {
	db = &mockDB{} // Use the mock so db is not nil

	body := []byte(`{"transactionId": 123, "amount": 999.99, "timestamp": `)
	req, err := http.NewRequest(http.MethodPost, "/api/transaction/", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createTransactionHandler)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("got %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestCreateTransactionHandler_ValidJSON(t *testing.T) {
	mock := &mockDB{}
	db = mock // Assign to global to avoid nil pointer

	tx := Transaction{
		TransactionID: "0f7e46df-c685-4df9-9e23-e75e7ac8ba7a",
		Amount:        99.99,
		Timestamp:     "2025-01-01T12:00:00Z",
	}
	bodyBytes, _ := json.Marshal(tx)

	req, err := http.NewRequest(http.MethodPost, "/api/transaction/", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createTransactionHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("got %v, want %v", rr.Code, http.StatusCreated)
	}

	expected := `{"message":"Transaction created successfully"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("got %v, want %v", rr.Body.String(), expected)
	}
}
