package registry

import (
	"fmt"
	"strings"

	"github.com/mohamedwael201193/pit/internal/config"
)

func RefuseMixedDesk(n config.Network, desk string) error {
	want := For(n).DeskID
	if desk == "" || want == "" {
		return nil
	}
	if !strings.EqualFold(desk, want) {
		return fmt.Errorf("desk_network_mix")
	}
	return nil
}
