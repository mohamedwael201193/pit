package companion

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChatWhatIsHappeningIdle(t *testing.T) {
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodPost, "/local/chat", bytes.NewBufferString(`{"text":"What's happening?"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Idle") {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), `"execute":true`) || strings.Contains(rec.Body.String(), `"posted":true`) {
		t.Fatal(rec.Body.String())
	}
}

func TestLocalModelsDirectOnly(t *testing.T) {
	h := New(t.TempDir())
	req := local(httptest.NewRequest(http.MethodGet, "/local/models", nil))
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Body.String())
	}
	body := strings.ToLower(rec.Body.String())
	if !strings.Contains(body, "glm-5.2") || !strings.Contains(body, `"path":"direct"`) {
		t.Fatal(rec.Body.String())
	}
	if strings.Contains(body, "router-api") || strings.Contains(body, "claude") {
		t.Fatal("router sku leaked")
	}
	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	models, _ := got["models"].([]any)
	if len(models) != 1 {
		t.Fatalf("catalog %d", len(models))
	}
}
