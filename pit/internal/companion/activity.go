package companion

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type activityEvent struct {
	TS          int64  `json:"ts"`
	WorkspaceID string `json:"workspace_id,omitempty"`
	Kind        string `json:"kind"`
	Market      string `json:"market,omitempty"`
	Action      string `json:"action,omitempty"`
	Status      string `json:"status,omitempty"`
	JobID       string `json:"job_id,omitempty"`
	PreviewHash string `json:"preview_hash,omitempty"`
	OID         string `json:"oid,omitempty"`
	Reason      string `json:"reason,omitempty"`
	Root        string `json:"root,omitempty"`
	Tx          string `json:"tx,omitempty"`
	TxLink      string `json:"tx_link,omitempty"`
	Digest      string `json:"digest,omitempty"`
	Link        string `json:"link,omitempty"`
	Autonomous  bool   `json:"autonomous,omitempty"`
	Sign        bool   `json:"sign"`
	Trade       bool   `json:"trade"`
}

func activityPath(dir string) string {
	return filepath.Join(dir, "activity.jsonl")
}

func appendActivity(dir string, ev activityEvent) {
	if dir == "" {
		return
	}
	ev.TS = time.Now().UnixMilli()
	ev.Sign = false
	ev.Trade = false
	raw, err := json.Marshal(ev)
	if err != nil || strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		return
	}
	f, err := os.OpenFile(activityPath(dir), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.Write(append(raw, '\n'))
}

func readActivity(dir string, limit int) []activityEvent {
	raw, err := os.ReadFile(activityPath(dir))
	if err != nil || strings.Contains(strings.ToLower(string(raw)), "app-sk-") {
		return nil
	}
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	out := make([]activityEvent, 0, len(lines))
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		var ev activityEvent
		if json.Unmarshal([]byte(ln), &ev) != nil {
			continue
		}
		out = append(out, ev)
	}
	if limit > 0 && len(out) > limit {
		out = out[len(out)-limit:]
	}
	return out
}
