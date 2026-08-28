package compute

import (
	"encoding/json"
	"os"

	"github.com/mohamedwael201193/pit/internal/engine"
)

func roleJSONFromOut(path string) []byte {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var wrap struct {
		Text string `json:"sanitized_output"`
	}
	if json.Unmarshal(b, &wrap) != nil {
		return nil
	}
	parsed := engine.ParseRoleText(wrap.Text)
	raw, err := json.Marshal(parsed)
	if err != nil {
		return nil
	}
	return raw
}
