package sdk

import (
	"fmt"
	"strings"

	"github.com/mohamedwael201193/pit/internal/config"
)

func (c Client) RefuseWrongExplorer(url string) error {
	want := config.For(c.Network).Explorer
	if !strings.HasPrefix(url, want) {
		return fmt.Errorf("explorer_mix")
	}
	return nil
}
