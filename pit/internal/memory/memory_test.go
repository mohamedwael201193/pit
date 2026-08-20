package memory

import (
	"strings"
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestKindsAndNoKeyInPrompt(t *testing.T) {
	if ValidKind("notes") {
		t.Fatal("kind")
	}
	ws := identity.NewWorkspaceID()
	key := "0x" + strings.Repeat("ab", 32)
	k, err := Put(config.Mainnet, ws, "observation", "book", []byte("x"), key)
	if err != nil {
		t.Fatal(err)
	}
	if err := PromptMayIncludeKey("analyze this", key); err != nil {
		t.Fatal(err)
	}
	if err := PromptMayIncludeKey("here is "+key, key); err == nil {
		t.Fatal("key leaked")
	}
	if k == "" {
		t.Fatal("key")
	}
}
