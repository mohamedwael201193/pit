package registry

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
)

func TestRefuseMixedDesk(t *testing.T) {
	main := For(config.Mainnet).DeskID
	if err := RefuseMixedDesk(config.Mainnet, main); err != nil {
		t.Fatal(err)
	}
	if main != "" {
		if err := RefuseMixedDesk(config.Mainnet, "0x0000000000000000000000000000000000000001"); err == nil {
			t.Fatal("mix")
		}
	}
}
