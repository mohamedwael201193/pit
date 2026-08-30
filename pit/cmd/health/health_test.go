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

func TestReleaseNeverSigns(t *testing.T) {
	fetchGH = func() (publicRelease, error) {
		return publicRelease{Tag: "v0.9.0", Name: "PIT 0.9.0", HTML: "https://github.com/mohamedwael201193/pit/releases/tag/v0.9.0", SHA: "ABCD", Unsigned: true}, nil
	}
	t.Cleanup(func() { fetchGH = fetchGitHubLatest })
	relMu.Lock()
	relOK = false
	relMu.Unlock()
	req := httptest.NewRequest(http.MethodGet, "/release", nil)
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
	if body["tag"] != "v0.9.0" {
		t.Fatal(body)
	}
}
