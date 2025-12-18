package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStatusHandler(t *testing.T) {
	h := NewHandler(nil) // Repo nie jest potrzebne dla StatusHandler
	req, err := http.NewRequest("GET", "/api/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.StatusHandler)

handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"service":"Fintech API","status":"OK","version":"2.0"}`
	if rr.Body.String() != expected+"\n" {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestTransactionHandler_InvalidMethod(t *testing.T) {
	h := NewHandler(nil)
	req, _ := http.NewRequest("GET", "/api/transactions", nil)
	rr := httptest.NewRecorder()

	h.TransactionHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rr.Code)
	}
}

func TestTransactionHandler_MalformedJSON(t *testing.T) {
	h := NewHandler(nil)
	req, _ := http.NewRequest("POST", "/api/transactions", strings.NewReader(`{invalid json}`))
	rr := httptest.NewRecorder()

	h.TransactionHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestAccountHandler_MissingUserID(t *testing.T) {
	h := NewHandler(nil)
	req, _ := http.NewRequest("GET", "/api/account/", nil)
	rr := httptest.NewRecorder()

	h.AccountHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}


