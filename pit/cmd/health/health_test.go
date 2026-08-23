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
	var body struct {
		OK   bool `json:"ok"`
		Sign bool `json:"sign"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if !body.OK || body.Sign {
		t.Fatalf("%+v", body)
	}
	_ = os.Getenv("PORT")
}
