package exec

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/hl"
)

func PostSigned(e *Exchange, env hl.Envelope, linked bool, previewHash, boundHash string) ([]byte, error) {
	if err := RefuseUnsigned(env.Signed()); err != nil {
		return nil, err
	}
	if err := RefusePostUntilLinked(linked); err != nil {
		return nil, err
	}
	if previewHash == "" || boundHash == "" || previewHash != boundHash {
		return nil, fmt.Errorf("preview_hash_mismatch")
	}
	return e.PostEnvelope(env)
}
