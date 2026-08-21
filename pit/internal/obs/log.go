package obs

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

type Event struct {
	RequestID string `json:"requestId"`
	Workspace string `json:"workspaceId,omitempty"`
	Phase     string `json:"phase"`
	Network   string `json:"network,omitempty"`
	Provider  string `json:"provider,omitempty"`
	Model     string `json:"model,omitempty"`
	Status    string `json:"status"`
	LatencyMs int64  `json:"latencyMs,omitempty"`
	Error     string `json:"error,omitempty"`
	At        string `json:"at"`
}

var forbidden = []string{
	"private_key", "session_key", "memory_key", "mnemonic", "seed",
}

func Sanitize(msg string) string {
	low := msg
	for _, w := range forbidden {
		if containsFold(low, w) {
			return "redacted"
		}
	}
	return msg
}

func containsFold(s, sub string) bool {
	if len(sub) == 0 || len(s) < len(sub) {
		return false
	}
	for i := 0; i+len(sub) <= len(s); i++ {
		ok := true
		for j := 0; j < len(sub); j++ {
			a, b := s[i+j], sub[j]
			if a >= 'A' && a <= 'Z' {
				a += 32
			}
			if b >= 'A' && b <= 'Z' {
				b += 32
			}
			if a != b {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}

func Write(w io.Writer, e Event) error {
	if w == nil {
		w = os.Stderr
	}
	e.At = time.Now().UTC().Format(time.RFC3339)
	e.Error = Sanitize(e.Error)
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	_, err = w.Write(append(b, '\n'))
	return err
}
