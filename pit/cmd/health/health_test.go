package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
)

func TestHealthNeverSigns(t *testing.T) {
	if err := config.GuardFallbacks(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["ok"] != true || body["sign"] != false {
		t.Fatalf("%+v", body)
	}
	if _, ok := body["wallet"]; ok {
		t.Fatal("wallet")
	}
}

func TestWatchNeverTrades(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/watch?network=mainnet", nil)
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["sign"] != false || body["trade"] != false {
		t.Fatalf("%+v", body)
	}
	if _, ok := body["private_book"]; ok {
		t.Fatal("book")
	}
	if _, ok := body["authorization"]; ok {
		t.Fatal("auth")
	}
}

func TestWatchCORSForProductionOrigin(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/watch?network=mainnet", nil)
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatal(rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "https://pit0g.vercel.app" {
		t.Fatal(rec.Header().Get("Access-Control-Allow-Origin"))
	}
	req = httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "https://evil.example")
	rec = httptest.NewRecorder()
	newMux().ServeHTTP(rec, req)
	if rec.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("evil origin")
	}
}
