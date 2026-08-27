package companion

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func local(r *http.Request) *http.Request {
	r.RemoteAddr = "127.0.0.1:9"
	return r
}

func TestPairingIssuesDeviceNotSession(t *testing.T) {
	h := New(t.TempDir())
	code, _ := h.Code()
	body, _ := json.Marshal(map[string]string{"code": code})
	req := local(httptest.NewRequest(http.MethodPost, "/pair", bytes.NewReader(body)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got["sign"] != false || got["canSign"] != false {
		t.Fatalf("%+v", got)
	}
	tok, _ := got["device"].(string)
	if len(tok) < 32 {
		t.Fatal("token")
	}
	if strings.Contains(rec.Body.String(), "private") {
		t.Fatal("secret")
	}

	req = local(httptest.NewRequest(http.MethodPost, "/pair", bytes.NewReader(body)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("replay")
	}
}

func TestMaliciousOriginCannotPair(t *testing.T) {
	h := New(t.TempDir())
	code, _ := h.Code()
	body, _ := json.Marshal(map[string]string{"code": code})
	req := local(httptest.NewRequest(http.MethodPost, "/pair", bytes.NewReader(body)))
	req.Header.Set("Origin", "https://evil.example")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal(rec.Code)
	}
}

func TestAuthorizeExportDenied(t *testing.T) {
	h := New(t.TempDir())
	for _, path := range []string{"/authorize", "/export", "/session"} {
		req := local(httptest.NewRequest(http.MethodPost, path, nil))
		req.Header.Set("Origin", "https://pit0g.vercel.app")
		rec := httptest.NewRecorder()
		h.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("%s %d", path, rec.Code)
		}
	}
}

func TestListenRejectsPublicBind(t *testing.T) {
	t.Setenv("PIT_COMPANION_ADDR", "0.0.0.0:17373")
	if _, err := ListenAddr(); err == nil {
		t.Fatal("public bind")
	}
}

func TestWrongCodeDenied(t *testing.T) {
	h := New(t.TempDir())
	body, _ := json.Marshal(map[string]string{"code": "AAAAAAAA"})
	req := local(httptest.NewRequest(http.MethodPost, "/pair", bytes.NewReader(body)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal(rec.Code)
	}
}

func TestLocalStatusAllowsLoopbackOmitsCode(t *testing.T) {
	h := New(t.TempDir())
	code, _ := h.Code()
	req := local(httptest.NewRequest(http.MethodGet, "/local/status", nil))
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), code) {
		t.Fatal("code leak")
	}
	req = local(httptest.NewRequest(http.MethodGet, "/local/status", nil))
	req.Header.Set("Origin", "https://evil.example")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal(rec.Code)
	}
}

func TestHealthOmitsPairingCode(t *testing.T) {
	h := New(t.TempDir())
	code, _ := h.Code()
	req := local(httptest.NewRequest(http.MethodGet, "/health", nil))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Code)
	}
	if strings.Contains(rec.Body.String(), code) {
		t.Fatal("code leak")
	}
}
