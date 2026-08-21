package exec

import (
	"encoding/json"
	"fmt"
)

func PostBound(e *Exchange, raw json.RawMessage, previewHash, boundHash string) ([]byte, error) {
	if previewHash == "" || boundHash == "" || previewHash != boundHash {
		return nil, fmt.Errorf("preview_hash_mismatch")
	}
	return e.Post(raw)
}
