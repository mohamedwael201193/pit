package exec

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func TestRefuseCooldown(t *testing.T) {
	p := policy.Default()
	p.CooldownSeconds = 60
	if err := RefuseCooldown(p, 100, 120); err == nil {
		t.Fatal("cool")
	}
	if err := RefuseCooldown(p, 100, 200); err != nil {
		t.Fatal(err)
	}
}
