package companion

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMalformedJSONDoesNotAuthorize(t *testing.T) {
	h := New(t.TempDir())
	payloads := []string{``, `{`, `[]`, `{"typed":"AUTHORIZE"}`, `{"hash":"0x1","typed":true}`, `null`}
	for _, p := range payloads {
		req := local(httptest.NewRequest(http.MethodPost, "/local/authorize", bytes.NewBufferString(p)))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		h.Handler().ServeHTTP(rec, req)
		if rec.Code == 200 && strings.Contains(rec.Body.String(), `"posted":true`) {
			t.Fatalf("posted %s %s", p, rec.Body.String())
		}
		if strings.Contains(strings.ToLower(rec.Body.String()), "app-sk-") {
			t.Fatal("secret")
		}
	}
}

func TestChatRefusesExecuteAndSecrets(t *testing.T) {
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodPost, "/local/chat", bytes.NewBufferString(`{"text":"buy ETH"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if !strings.Contains(rec.Body.String(), `"execute":false`) {
		t.Fatal(rec.Body.String())
	}
	req = local(httptest.NewRequest(http.MethodPost, "/local/chat", bytes.NewBufferString(`{"text":"Bearer app-sk-secret"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if !strings.Contains(rec.Body.String(), "secret_refused") {
		t.Fatal(rec.Body.String())
	}
}

func TestChatResearchDoesNotExecute(t *testing.T) {
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodPost, "/local/chat", bytes.NewBufferString(`{"text":"Research ETH privately."}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if strings.Contains(rec.Body.String(), `"execute":true`) || strings.Contains(rec.Body.String(), `"posted":true`) {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"start_research":true`) {
		t.Fatal(rec.Body.String())
	}
}

func TestChatGreetingIsNotCannedHelp(t *testing.T) {
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodPost, "/local/chat", bytes.NewBufferString(`{"text":"HI"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if strings.Contains(rec.Body.String(), "I can research a policy market") {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Hello") {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), `"execute":true`) {
		t.Fatal(rec.Body.String())
	}
}

func TestChatCannotExecuteDecorated(t *testing.T) {
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodPost, "/local/chat", bytes.NewBufferString(`{"text":"Why can't PIT execute this?"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), `"execute":true`) || strings.Contains(rec.Body.String(), `"posted":true`) {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "AUTHORIZE") {
		t.Fatal(rec.Body.String())
	}
}

func TestChatTypoResearchDoesNotExecute(t *testing.T) {
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodPost, "/local/chat", bytes.NewBufferString(`{"text":"HI WHAT PRICE OF ETH NOW AND GIVE ME REASEARCH AND DO TRADE"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if strings.Contains(rec.Body.String(), `"execute":true`) || strings.Contains(rec.Body.String(), `"posted":true`) {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"start_research":true`) {
		t.Fatal(rec.Body.String())
	}
}

func TestAutomationCannotAuthorize(t *testing.T) {
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodPost, "/local/automation", bytes.NewBufferString(`{"watch":true,"auto_research":true,"execute":true}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if strings.Contains(rec.Body.String(), `"execute":true`) {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "automation_cannot_authorize") {
		t.Fatal(rec.Body.String())
	}
}
