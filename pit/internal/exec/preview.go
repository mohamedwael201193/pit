package exec

import (
	"strings"

	"github.com/mohamedwael201193/pit/internal/engine"
)

func HashForAuthorize(display string) string {
	return strings.TrimPrefix(display, "0x")
}

func RequirePreview(p engine.Preview, displayHash string, nowMs int64) error {
	return engine.Authorize(p, HashForAuthorize(displayHash), nowMs, map[string]struct{}{})
}
