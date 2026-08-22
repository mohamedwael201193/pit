package exec

import (
	"encoding/json"
	"fmt"

	"github.com/mohamedwael201193/pit/internal/hl"
)

func CancelBound(e *Exchange, asset int, cloid, previewHash, boundHash string) ([]byte, error) {
	if previewHash == "" || boundHash == "" || previewHash != boundHash {
		return nil, fmt.Errorf("preview_hash_mismatch")
	}
	raw, err := hl.BuildCancel(asset, cloid)
	if err != nil {
		return nil, err
	}
	return PostBound(e, json.RawMessage(raw), previewHash, boundHash)
}
