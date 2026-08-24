package watch

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func TestKillBlocks(t *testing.T) {
	p := policy.Default()
	if err := KillBlocks(p); err != nil {
		t.Fatal(err)
	}
	p.KillSwitch = true
	if err := KillBlocks(p); err == nil {
		t.Fatal("kill")
	}
}
