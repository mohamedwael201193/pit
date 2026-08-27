package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthAllowsProductionOrigin(t *testing.T) {
	if !HealthOriginOK("https://pit0g.vercel.app") {
		t.Fatal("prod")
	}
	if HealthOriginOK("https://evil.example") {
		t.Fatal("evil")
	}
}

func TestCompanionRejectsForeignOrigin(t *testing.T) {
	if CompanionOriginOK("https://evil.example") {
		t.Fatal("evil")
	}
	if !CompanionOriginOK("https://pit0g.vercel.app") {
		t.Fatal("prod")
	}
}

func TestCompanionLoopbackAndOrigin(t *testing.T) {
	ok := Companion(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.RemoteAddr = "8.8.8.8:9"
	rec := httptest.NewRecorder()
	ok.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("remote %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/health", nil)
	req.RemoteAddr = "127.0.0.1:9"
	req.Header.Set("Origin", "https://evil.example")
	rec = httptest.NewRecorder()
	ok.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("origin %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/health", nil)
	req.RemoteAddr = "127.0.0.1:9"
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec = httptest.NewRecorder()
	ok.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("allow %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "https://pit0g.vercel.app" {
		t.Fatal(rec.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestPublicCORSPreflight(t *testing.T) {
	h := Public(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodOptions, "/watch", nil)
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatal(rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "https://pit0g.vercel.app" {
		t.Fatal("acao")
	}
}
