package companion

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChatThreadsIsolatedAndRename(t *testing.T) {
	dir := t.TempDir()
	h := New(dir)
	req := local(httptest.NewRequest(http.MethodPost, "/local/chat/thread", bytes.NewBufferString(`{"op":"new","title":"New desk"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	var created map[string]any
	if json.Unmarshal(rec.Body.Bytes(), &created) != nil || created["id"] == nil {
		t.Fatal(rec.Body.String())
	}
	id := created["id"].(string)
	req = local(httptest.NewRequest(http.MethodPost, "/local/chat", bytes.NewBufferString(`{"text":"Why is ETH interesting?","thread":"`+id+`"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	req = local(httptest.NewRequest(http.MethodGet, "/local/chat/log?thread="+id, nil))
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if !bytes.Contains(rec.Body.Bytes(), []byte("ETH")) {
		t.Fatal(rec.Body.String())
	}
	req = local(httptest.NewRequest(http.MethodGet, "/local/chat/log?thread=desk", nil))
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if bytes.Contains(rec.Body.Bytes(), []byte("ETH")) && bytes.Contains(rec.Body.Bytes(), []byte(id)) {
		t.Fatal("desk saw other thread id in messages")
	}
	req = local(httptest.NewRequest(http.MethodPost, "/local/chat/thread", bytes.NewBufferString(`{"op":"delete","id":"desk"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if bytes.Contains(rec.Body.Bytes(), []byte(`"ok":true`)) && !bytes.Contains(rec.Body.Bytes(), []byte("thread_protected")) {
		t.Fatal("desk thread deleted")
	}
}
