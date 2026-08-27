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

func TestLocalCodeDesktopOnly(t *testing.T) {
	h := New(t.TempDir())
	want, _ := h.Code()

	req := local(httptest.NewRequest(http.MethodGet, "/local/code", nil))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("web must not read pairing code")
	}

	req = local(httptest.NewRequest(http.MethodGet, "/local/code", nil))
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), want) {
		t.Fatal("desktop code")
	}

	req = local(httptest.NewRequest(http.MethodGet, "/local/code", nil))
	req.Header.Set("Origin", "http://tauri.localhost")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal("windows webview pairing")
	}

	req = local(httptest.NewRequest(http.MethodGet, "/local/doctor", nil))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("web doctor")
	}
}

func TestLocalStatusVersionNoSecret(t *testing.T) {
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodGet, "/local/status", nil))
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got["sign"] == true || got["trade"] == true {
		t.Fatal(got)
	}
	if got["version"] != "0.1.4" {
		t.Fatalf("version %v", got["version"])
	}
}

func TestDesktopInitSessionPolicyUserB(t *testing.T) {
	h := New(t.TempDir())
	a := "0x1111111111111111111111111111111111111111"
	b := "0x2222222222222222222222222222222222222222"

	req := local(httptest.NewRequest(http.MethodPost, "/local/init", strings.NewReader(`{"wallet":"`+a+`","network":"mainnet"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(strings.ToLower(rec.Body.String()), "private") {
		t.Fatal("secret")
	}

	req = local(httptest.NewRequest(http.MethodPost, "/local/init", strings.NewReader(`{"wallet":"`+b+`","network":"mainnet"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("user B %d %s", rec.Code, rec.Body.String())
	}

	req = local(httptest.NewRequest(http.MethodPost, "/local/init", strings.NewReader(`{"wallet":"`+a+`","network":"mainnet"}`)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("web must not init")
	}

	req = local(httptest.NewRequest(http.MethodPost, "/local/session", nil))
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	var sess map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &sess); err != nil {
		t.Fatal(err)
	}
	if sess["withdraw"] != false || sess["sign"] != false || sess["agent"] == nil || sess["agent"] == "" {
		t.Fatal(sess)
	}
	if strings.Contains(strings.ToLower(rec.Body.String()), "private_key") {
		t.Fatal("key leak")
	}

	req = local(httptest.NewRequest(http.MethodPost, "/local/session", nil))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("web session")
	}

	req = local(httptest.NewRequest(http.MethodPost, "/local/policy", nil))
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}

	req = local(httptest.NewRequest(http.MethodGet, "/local/status", nil))
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), a) {
		t.Fatal("wallet missing")
	}
}

func TestBrowserBindRequiresDeviceAndIsolatesUserB(t *testing.T) {
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
	var paired map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &paired); err != nil {
		t.Fatal(err)
	}
	tok, _ := paired["device"].(string)

	req = local(httptest.NewRequest(http.MethodPost, "/bind", strings.NewReader(`{"wallet":"0x1111111111111111111111111111111111111111","network":"mainnet"}`)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatal("no bearer")
	}

	req = local(httptest.NewRequest(http.MethodPost, "/bind", strings.NewReader(`{"wallet":"0x1111111111111111111111111111111111111111","network":"mainnet"}`)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tok)
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}

	req = local(httptest.NewRequest(http.MethodPost, "/bind", strings.NewReader(`{"wallet":"0x2222222222222222222222222222222222222222","network":"mainnet"}`)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tok)
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("user B bind %d %s", rec.Code, rec.Body.String())
	}
}

func TestDirectIntentNeverLeaksToken(t *testing.T) {
	t.Setenv("PIT_KEYRING", "file")
	t.Setenv("PIT_DIRECT_AUTH_FILE", "")
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodPost, "/local/init", strings.NewReader(`{"wallet":"0x1111111111111111111111111111111111111111","network":"mainnet"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}

	req = local(httptest.NewRequest(http.MethodGet, "/local/direct-intent", nil))
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	body := rec.Body.String()
	if strings.Contains(strings.ToLower(body), "app-sk-") {
		t.Fatal("token leak")
	}
	if !strings.Contains(body, `"digest"`) || !strings.Contains(body, `"message"`) {
		t.Fatal(body)
	}

	req = local(httptest.NewRequest(http.MethodGet, "/local/direct-intent", nil))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("web must not read local direct intent")
	}

	req = local(httptest.NewRequest(http.MethodPost, "/direct/intent", nil))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatal("device required")
	}

	req = local(httptest.NewRequest(http.MethodPost, "/local/research", strings.NewReader(`{"coin":"ETH"}`)))
	req.Header.Set("Origin", "https://pit0g.vercel.app")
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatal("web research")
	}
}
