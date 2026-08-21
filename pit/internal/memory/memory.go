package memory

import (
	"fmt"
	"os"
	"strings"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/storage"
)

var Kinds = []string{
	"observation", "thesis", "forecast", "execution", "outcome", "error",
	"calibration", "role_performance", "skill_performance", "risk_event", "policy_decision", "regime",
}

func ValidKind(k string) bool {
	for _, x := range Kinds {
		if x == k {
			return true
		}
	}
	return false
}

func Put(net config.Network, workspaceID, kind, name string, payload []byte, keyHex string) (string, error) {
	if !ValidKind(kind) {
		return "", fmt.Errorf("bad_kind")
	}
	if err := storage.RequireHexKey(keyHex); err != nil {
		return "", err
	}
	if os.Getenv("PIT_MEMORY_KEY") != "" && os.Getenv("PIT_PRODUCT_MODE") == "true" {
		return "", fmt.Errorf("global_memory_key_forbidden")
	}
	return storage.ObjectKey(net, workspaceID, kind, name)
}

func PromptMayIncludeKey(prompt string, keyHex string) error {
	if keyHex != "" && strings.Contains(prompt, keyHex) {
		return fmt.Errorf("memory_key_in_prompt")
	}
	return nil
}
