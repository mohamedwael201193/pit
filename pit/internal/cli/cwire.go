package cli

import (
	"encoding/json"

	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/session"
)

func CancelWire(asset int, cloid string) (json.RawMessage, error) {
	if err := session.CheckAction("cancel"); err != nil {
		return nil, err
	}
	return hl.BuildCancel(asset, cloid)
}
