package companion

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mohamedwael201193/pit/internal/httpx"
)

type chatThread struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Created int64  `json:"created"`
	Updated int64  `json:"updated"`
	Preview string `json:"preview"`
}

func threadsPath(dir string) string { return filepath.Join(dir, "chat-threads.enc") }

func defaultThread() chatThread {
	now := time.Now().UnixMilli()
	return chatThread{ID: "desk", Title: "Desk", Created: now, Updated: now, Preview: ""}
}

func loadThreads(dir string) []chatThread {
	raw, err := os.ReadFile(threadsPath(dir))
	if err != nil {
		return []chatThread{defaultThread()}
	}
	plain, err := openBytes(dir, string(raw))
	if err != nil {
		return []chatThread{defaultThread()}
	}
	var rows []chatThread
	if json.Unmarshal(plain, &rows) != nil || len(rows) == 0 {
		return []chatThread{defaultThread()}
	}
	return rows
}

func saveThreads(dir string, rows []chatThread) {
	if len(rows) == 0 {
		rows = []chatThread{defaultThread()}
	}
	raw, err := json.Marshal(rows)
	if err != nil {
		return
	}
	sealed, err := sealBytes(dir, raw)
	if err != nil {
		return
	}
	_ = os.WriteFile(threadsPath(dir), []byte(sealed), 0o600)
}

func newThreadID() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func touchThread(dir, id, preview string) {
	if strings.TrimSpace(id) == "" {
		id = "desk"
	}
	rows := loadThreads(dir)
	found := false
	now := time.Now().UnixMilli()
	clip := preview
	if len(clip) > 80 {
		clip = clip[:80]
	}
	for i := range rows {
		if rows[i].ID == id {
			rows[i].Updated = now
			if clip != "" {
				rows[i].Preview = clip
			}
			if rows[i].Title == "New desk" && clip != "" {
				t := clip
				if len(t) > 40 {
					t = t[:40]
				}
				rows[i].Title = t
			}
			found = true
		}
	}
	if !found {
		rows = append(rows, chatThread{ID: id, Title: "Desk", Created: now, Updated: now, Preview: clip})
	}
	saveThreads(dir, rows)
}

func (h *Hub) localChatThreads(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !httpx.CodeOriginOK(r.Header.Get("Origin")) {
		http.Error(w, "origin_denied", http.StatusForbidden)
		return
	}
	writeLocal(w, http.StatusOK, map[string]any{"threads": loadThreads(h.Dir), "sign": false, "trade": false})
}

func (h *Hub) localChatThread(w http.ResponseWriter, r *http.Request) {
	if !desktopOnly(w, r) {
		return
	}
	var body struct {
		Op    string `json:"op"`
		ID    string `json:"id"`
		Title string `json:"title"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	rows := loadThreads(h.Dir)
	switch strings.ToLower(strings.TrimSpace(body.Op)) {
	case "new":
		now := time.Now().UnixMilli()
		id := newThreadID()
		title := strings.TrimSpace(body.Title)
		if title == "" {
			title = "New desk"
		}
		rows = append([]chatThread{{ID: id, Title: title, Created: now, Updated: now}}, rows...)
		saveThreads(h.Dir, rows)
		writeLocal(w, http.StatusOK, map[string]any{"ok": true, "id": id, "threads": rows, "sign": false, "trade": false})
		return
	case "rename":
		title := strings.TrimSpace(body.Title)
		if body.ID == "" || title == "" {
			writeLocal(w, http.StatusOK, map[string]any{"ok": false, "error": "thread_required", "sign": false, "trade": false})
			return
		}
		for i := range rows {
			if rows[i].ID == body.ID {
				rows[i].Title = title
				rows[i].Updated = time.Now().UnixMilli()
			}
		}
		saveThreads(h.Dir, rows)
		writeLocal(w, http.StatusOK, map[string]any{"ok": true, "id": body.ID, "threads": rows, "sign": false, "trade": false})
		return
	case "delete":
		if body.ID == "" || body.ID == "desk" {
			writeLocal(w, http.StatusOK, map[string]any{"ok": false, "error": "thread_protected", "threads": rows, "sign": false, "trade": false})
			return
		}
		out := make([]chatThread, 0, len(rows))
		for _, t := range rows {
			if t.ID != body.ID {
				out = append(out, t)
			}
		}
		if len(out) == 0 {
			out = []chatThread{defaultThread()}
		}
		saveThreads(h.Dir, out)
		writeLocal(w, http.StatusOK, map[string]any{"ok": true, "threads": out, "sign": false, "trade": false})
		return
	default:
		writeLocal(w, http.StatusOK, map[string]any{"ok": false, "error": "unknown_op", "sign": false, "trade": false})
	}
}
