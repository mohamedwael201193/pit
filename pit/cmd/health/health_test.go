package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
)

func TestHealthNeverSigns(t *testing.T) {
	if err := config.GuardFallbacks(); err != nil {
		t.Fatal(err)
	}
	if err := config.RefuseSessionEnv(); err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "service": "pit", "sign": false})
	})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
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
	if _, ok := body["session"]; ok {
		t.Fatal("session")
	}
	_ = os.Getenv("PORT")
}
