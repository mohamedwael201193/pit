package obs

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mohamedwael201193/pit/internal/version"
	"github.com/mohamedwael201193/pit/internal/watch"
)

func WriteJSON(w http.ResponseWriter, status int, body map[string]any) {
	if err := RefuseHealthSecrets(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func HealthBody(requestID string) map[string]any {
	return map[string]any{
		"ok":        true,
		"service":   "pit",
		"time":      time.Now().UTC().Format(time.RFC3339),
		"sign":      false,
		"requestId": requestID,
		"version":   version.Number,
	}
}

func WatchBody(v watch.PublicView) map[string]any {
	b, _ := json.Marshal(v)
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	if m == nil {
		m = map[string]any{}
	}
	m["sign"] = false
	m["trade"] = false
	return m
}
