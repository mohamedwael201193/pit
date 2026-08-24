package exec

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/engine"
)

func RefuseStalePreview(p engine.Preview, hash string, nowMs int64) error {
	used := map[string]struct{}{}
	if err := engine.Authorize(p, hash, nowMs, used); err != nil {
		return err
	}
	if nowMs >= p.ExpiryUnixMs {
		return fmt.Errorf("preview_expired")
	}
	return nil
}
