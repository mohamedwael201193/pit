package companion

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorkspaceFilesDoNotCross(t *testing.T) {
	a := t.TempDir()
	b := t.TempDir()
	appendActivity(a, activityEvent{WorkspaceID: "wa", Kind: "research.started", JobID: "job-a"})
	appendChat(a, "user", "secret thesis for A", "")
	if len(readActivity(b, 20)) != 0 {
		t.Fatal("B saw A activity")
	}
	if len(readChat(b, 20)) != 0 {
		t.Fatal("B saw A chat")
	}
	ha := New(a)
	hb := New(b)
	req := local(httptest.NewRequest(http.MethodGet, "/local/activity", nil))
	rec := httptest.NewRecorder()
	hb.Handler().ServeHTTP(rec, req)
	if strings.Contains(rec.Body.String(), "job-a") {
		t.Fatal(rec.Body.String())
	}
	req = local(httptest.NewRequest(http.MethodGet, "/local/chat/log", nil))
	rec = httptest.NewRecorder()
	ha.Handler().ServeHTTP(rec, req)
	if strings.Contains(rec.Body.String(), "app-sk-") {
		t.Fatal("secret")
	}
}

func TestChatCannotPostOrder(t *testing.T) {
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodPost, "/local/chat", bytes.NewBufferString(`{"text":"I authorize it"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if strings.Contains(rec.Body.String(), `"posted":true`) || strings.Contains(rec.Body.String(), `"execute":true`) {
		t.Fatal(rec.Body.String())
	}
}

func TestForeignOriginDeniedOnSurfaces(t *testing.T) {
	h := New(t.TempDir())
	req := httptest.NewRequest(http.MethodGet, "/local/activity", nil)
	req.Header.Set("Origin", "https://evil.example")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code == 200 && strings.Contains(rec.Body.String(), `"events"`) {
		t.Fatal(rec.Body.String())
	}
}

func TestForgetDoesNotDeleteReceipts(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "last-order.json"), []byte(`{"oid":"1","sign":false}`), 0o600)
	appendChat(dir, "user", "hello", "")
	h := New(dir)
	req := local(httptest.NewRequest(http.MethodPost, "/local/memory/forget", bytes.NewBufferString(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if _, err := os.Stat(filepath.Join(dir, "last-order.json")); err != nil {
		t.Fatal("receipt wiped")
	}
}
